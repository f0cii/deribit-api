package multicast

import (
	"errors"
	"fmt"
	"io"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/KyberNetwork/deribit-api/pkg/multicast/sbe"
	"github.com/shopspring/decimal"
)

type EventType int

// Defines constants for types of events.
const (
	EventTypeInstrument EventType = iota
	EventTypeOrderBook
	EventTypeTrades
	EventTypeTicker
)

// Event represents a Deribit multicast events.
type Event struct {
	Type EventType
	Data interface{}
}

// DecodeEvents decodes a UDP package into a list of events.
func DecodeEvents(m *sbe.SbeGoMarshaller, r io.Reader, instrumentIDToName map[uint32]string) (events []Event, err error) {
	// Decode Sbe messages
	for {
		// Decode message header.
		var header sbe.MessageHeader
		err = header.Decode(m, r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}
		event, err := DecodeEvent(m, r, header, instrumentIDToName)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
}

func DecodeEvent(
	m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader,
	instrumentIDToName map[uint32]string,
) (Event, error) {
	switch header.TemplateId {
	case 1000:
		return DecodeInstrumentEvent(m, r, header)
	case 1001:
		return DecodeOrderBookEvent(m, r, header)
	case 1002:
		return DecodeTradesEvent(m, r, header, instrumentIDToName)
	case 1003:
		return DecodeTickerEvent(m, r, header, instrumentIDToName)
	default:
		return Event{}, nil
	}
}

func DecodeInstrumentEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var ins sbe.Instrument
	err := ins.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	instrument := models.Instrument{
		TickSize:             decimal.NewFromFloat(ins.TickSize),
		TakerCommission:      decimal.NewFromFloat(ins.TakerCommission),
		SettlementPeriod:     ins.SettlementPeriod.String(),
		QuoteCurrency:        string(ins.QuoteCurrency[:]),
		MinTradeAmount:       decimal.NewFromFloat(ins.MinTradeAmount),
		MakerCommission:      decimal.NewFromFloat(ins.MakerCommission),
		Leverage:             int(ins.MaxLeverage),
		Kind:                 ins.Kind.String(),
		IsActive:             ins.InstrumentState.IsActive(),
		InstrumentID:         ins.InstrumentId,
		InstrumentName:       string(ins.InstrumentName),
		ExpirationTimestamp:  ins.ExpirationTimestampMs,
		CreationTimestamp:    ins.CreationTimestampMs,
		ContractSize:         decimal.NewFromFloat(ins.ContractSize),
		BaseCurrency:         string(ins.BaseCurrency[:]),
		BlockTradeCommission: decimal.NewFromFloat(ins.BlockTradeCommission),
		OptionType:           ins.OptionType.String(),
		Strike:               decimal.NewFromFloat(ins.StrikePrice),
	}
	return Event{
		Type: EventTypeInstrument,
		Data: instrument,
	}, nil
}

func DecodeOrderBookEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var b sbe.Book
	err := b.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	book := models.OrderBookRawNotification{
		Timestamp:      int64(b.TimestampMs),
		InstrumentName: fmt.Sprintf("%v", b.InstrumentId),
		PrevChangeID:   int64(b.PrevChangeId),
		ChangeID:       int64(b.ChangeId),
	}

	for _, bookChange := range b.ChangesList {
		item := models.OrderBookNotificationItem{
			Action: bookChange.Change.String(),
			Price:  decimal.NewFromFloat(bookChange.Price),
			Amount: decimal.NewFromFloat(bookChange.Amount),
		}

		if bookChange.Side == sbe.BookSide.Ask {
			book.Asks = append(book.Asks, item)
		} else if bookChange.Side == sbe.BookSide.Bid {
			book.Bids = append(book.Bids, item)
		}
	}

	return Event{
		Type: EventTypeOrderBook,
		Data: book,
	}, nil
}

func DecodeTradesEvent(
	m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader,
	instrumentIDToName map[uint32]string,
) (Event, error) {
	var t sbe.Trades
	err := t.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	trades := make(models.TradesNotification, len(t.TradesList))
	for id, trade := range t.TradesList {
		trades[id] = models.Trade{
			Amount:         decimal.NewFromFloat(trade.Amount),
			BlockTradeID:   fmt.Sprintf("%d", trade.BlockTradeId),
			Direction:      trade.Direction.String(),
			IndexPrice:     decimal.NewFromFloat(trade.IndexPrice),
			InstrumentName: instrumentIDToName[t.InstrumentId],
			IV:             decimal.NewFromFloat(trade.Iv),
			Liquidation:    trade.Liquidation.String(),
			MarkPrice:      decimal.NewFromFloat(trade.MarkPrice),
			Price:          decimal.NewFromFloat(trade.Price),
			TickDirection:  int(trade.TickDirection),
			Timestamp:      trade.TimestampMs,
			TradeID:        fmt.Sprintf("%d", trade.TradeId),
			TradeSeq:       trade.TradeSeq,
		}
	}

	return Event{
		Type: EventTypeTrades,
		Data: trades,
	}, nil
}

func DecodeTickerEvent(
	m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader,
	instrumentIDToName map[uint32]string,
) (Event, error) {
	var t sbe.Ticker
	err := t.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	bestBidPrice := decimal.NewFromFloat(t.BestBidPrice)
	bestAskPrice := decimal.NewFromFloat(t.BestAskPrice)

	ticker := models.TickerNotification{
		Timestamp:       t.TimestampMs,
		Stats:           models.Stats{},
		State:           t.InstrumentState.String(),
		SettlementPrice: decimal.NewFromFloat(t.SettlementPrice),
		OpenInterest:    decimal.NewFromFloat(t.OpenInterest),
		MinPrice:        decimal.NewFromFloat(t.MinSellPrice),
		MaxPrice:        decimal.NewFromFloat(t.MaxBuyPrice),
		MarkPrice:       decimal.NewFromFloat(t.MarkPrice),
		LastPrice:       decimal.NewFromFloat(t.LastPrice),
		InstrumentName:  instrumentIDToName[t.InstrumentId],
		IndexPrice:      decimal.NewFromFloat(t.IndexPrice),
		Funding8H:       decimal.NewFromFloat(t.Funding8h),
		CurrentFunding:  decimal.NewFromFloat(t.CurrentFunding),
		BestBidPrice:    &bestBidPrice,
		BestBidAmount:   decimal.NewFromFloat(t.BestBidAmount),
		BestAskPrice:    &bestAskPrice,
		BestAskAmount:   decimal.NewFromFloat(t.BestAskAmount),
	}
	return Event{
		Type: EventTypeTicker,
		Data: ticker,
	}, nil
}
