package multicast

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"sync"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/KyberNetwork/deribit-api/pkg/multicast/sbe"
	"github.com/KyberNetwork/deribit-api/pkg/websocket"
	"github.com/chuckpreslar/emission"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const (
	defaultReadBufferSize = 1500
	defaultDataChSize     = 20
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
	Instrument string
	Type       EventType
	Data       interface{}
}

// Client represents a client for Deribit multicast.
type Client struct {
	m           *sbe.SbeGoMarshaller
	log         *zap.SugaredLogger
	mu          sync.RWMutex
	inf         *net.Interface
	addrs       []string
	addrConnMap map[string]*net.UDPConn
	wsClient    *websocket.Client

	isConnected        bool
	supportCurrencies  []string
	instrumentIDToName map[uint32]string
	channelIDSeqNum    map[uint16]uint32
	emitter            *emission.Emitter
}

// NewClient creates a new Client instance.
func NewClient(ifname string, addrs []string, wsClient *websocket.Client, currencies []string) (*Client, error) {
	l := zap.S()

	inf, err := net.InterfaceByName(ifname)
	if err != nil {
		l.Errorw("failed to create net interfaces by name", "err", err)
		return nil, err
	}

	client := &Client{
		m:           sbe.NewSbeGoMarshaller(),
		log:         l,
		mu:          sync.RWMutex{},
		inf:         inf,
		addrs:       addrs,
		wsClient:    wsClient,
		addrConnMap: make(map[string]*net.UDPConn),

		supportCurrencies:  currencies,
		instrumentIDToName: make(map[uint32]string),
		channelIDSeqNum:    make(map[uint16]uint32),
		emitter:            emission.NewEmitter(),
	}

	err = client.buildInstrumentIDToNameMapping()
	if err != nil {
		l.Errorw("failed to build instrumentId name mapping", "err", err)
		return nil, err
	}

	err = client.Start()
	if err != nil {
		client.log.Errorw("Fail to start multicast connection", "error", err)
		return nil, err
	}

	return client, nil
}

// BuildInstrumentIDToNameMapping builds a mapping to map instrument id to instrument name.
func (c *Client) buildInstrumentIDToNameMapping() error {
	// need to update this field via instruments_update event
	instruments, err := GetAllInstrument(c.wsClient, c.supportCurrencies)
	if err != nil {
		c.log.Errorw("failed to get all instruments", "err", err)
		return err
	}
	for _, instrument := range instruments {
		c.instrumentIDToName[instrument.InstrumentID] = instrument.InstrumentName
	}
	return nil
}

// GetAllInstrument returns a list of all instruments by currencies
func GetAllInstrument(wsClient *websocket.Client, currencies []string) ([]models.Instrument, error) {
	result := make([]models.Instrument, 0)
	for _, currency := range currencies {
		ins, err := wsClient.GetInstruments(
			context.Background(),
			&models.GetInstrumentsParams{
				Currency: currency,
			})
		if err != nil {
			return nil, err
		}
		result = append(result, ins...)
	}
	return result, nil
}

// DecodeEvents decodes a UDP package into a list of events.
func (c *Client) decodeEvents(r io.Reader) (events []Event, err error) {
	// Decode Sbe messages
	for {
		// Decode message header.
		var header sbe.MessageHeader
		err = header.Decode(c.m, r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}
		event, err := c.decodeEvent(r, header)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
}

func (c *Client) decodeEvent(r io.Reader, header sbe.MessageHeader) (Event, error) {
	switch header.TemplateId {
	case 1000:
		return c.decodeInstrumentEvent(r, header)
	case 1001:
		return c.decodeOrderBookEvent(r, header)
	case 1002:
		return c.decodeTradesEvent(r, header)
	case 1003:
		return c.decodeTickerEvent(r, header)
	default:
		return Event{}, nil
	}
}

func (c *Client) decodeInstrumentEvent(r io.Reader, header sbe.MessageHeader) (Event, error) {
	var ins sbe.Instrument
	err := ins.Decode(c.m, r, header.BlockLength, true)
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

func (c *Client) decodeOrderBookEvent(r io.Reader, header sbe.MessageHeader) (Event, error) {
	var b sbe.Book
	err := b.Decode(c.m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	c.mu.Lock()
	instrumentName := c.instrumentIDToName[b.InstrumentId]
	c.mu.Unlock()

	book := models.OrderBookRawNotification{
		Timestamp:      int64(b.TimestampMs),
		InstrumentName: instrumentName,
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
		Instrument: instrumentName,
		Type:       EventTypeOrderBook,
		Data:       book,
	}, nil
}

func (c *Client) decodeTradesEvent(r io.Reader, header sbe.MessageHeader) (Event, error) {
	var t sbe.Trades
	err := t.Decode(c.m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	c.mu.Lock()
	instrumentName := c.instrumentIDToName[t.InstrumentId]
	c.mu.Unlock()

	trades := make(models.TradesNotification, len(t.TradesList))
	for id, trade := range t.TradesList {

		trades[id] = models.Trade{
			Amount:         decimal.NewFromFloat(trade.Amount),
			BlockTradeID:   strconv.FormatUint(trade.BlockTradeId, 10),
			Direction:      trade.Direction.String(),
			IndexPrice:     decimal.NewFromFloat(trade.IndexPrice),
			InstrumentName: instrumentName,
			IV:             decimal.NewFromFloat(trade.Iv),
			Liquidation:    trade.Liquidation.String(),
			MarkPrice:      decimal.NewFromFloat(trade.MarkPrice),
			Price:          decimal.NewFromFloat(trade.Price),
			TickDirection:  int(trade.TickDirection),
			Timestamp:      trade.TimestampMs,
			TradeID:        strconv.FormatUint(trade.BlockTradeId, 10),
			TradeSeq:       trade.TradeSeq,
		}
	}

	return Event{
		Instrument: instrumentName,
		Type:       EventTypeTrades,
		Data:       trades,
	}, nil
}

func (c *Client) decodeTickerEvent(r io.Reader, header sbe.MessageHeader) (Event, error) {
	var t sbe.Ticker
	err := t.Decode(c.m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	bestBidPrice := decimal.NewFromFloat(t.BestBidPrice)
	bestAskPrice := decimal.NewFromFloat(t.BestAskPrice)

	c.mu.Lock()
	instrumentName := c.instrumentIDToName[t.InstrumentId]
	c.mu.Unlock()

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
		InstrumentName:  instrumentName,
		IndexPrice:      decimal.NewFromFloat(t.IndexPrice),
		Funding8H:       decimal.NewFromFloat(t.Funding8h),
		CurrentFunding:  decimal.NewFromFloat(t.CurrentFunding),
		BestBidPrice:    &bestBidPrice,
		BestBidAmount:   decimal.NewFromFloat(t.BestBidAmount),
		BestAskPrice:    &bestAskPrice,
		BestAskAmount:   decimal.NewFromFloat(t.BestAskAmount),
	}
	return Event{
		Instrument: instrumentName,
		Type:       EventTypeTicker,
		Data:       ticker,
	}, nil
}

// IsConnected checks whether the connection is established or not.
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}

// Start starts listening events on interface `ifname` and `addrs`.
func (c *Client) Start() error {
	// initiate connections
	for _, addr := range c.addrs {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			c.log.Errorw("failed to resolve UDP address", "err", err)
			return err
		}

		udpConn, err := net.ListenMulticastUDP("udp", c.inf, udpAddr)
		if err != nil {
			c.log.Errorw("failed to listen to multicast UDP", "error", err)
			return err
		}

		err = udpConn.SetReadBuffer(defaultReadBufferSize)
		if err != nil {
			c.log.Errorw("fail to set read buffer for UDP", "error", err)
			return err
		}

		c.addrConnMap[addr] = udpConn
	}

	return c.ListenToEvents()
}

// Stop stops listening for events.
func (c *Client) Stop() error {
	for _, conn := range c.addrConnMap {
		err := conn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getSeqNum(channelId uint16) (seq uint32, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	seq, ok = c.channelIDSeqNum[channelId]
	return
}

func (c *Client) setSeqNum(channelId uint16, seq uint32) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.channelIDSeqNum[channelId] = seq
}

func readPackageHeader(r io.Reader) (uint16, uint16, uint32, error) {
	b8 := make([]byte, 8)
	if _, err := io.ReadFull(r, b8); err != nil {
		return 0, 0, 0, err
	}

	n := uint16(b8[0]) | uint16(b8[1])<<8
	channelID := uint16(b8[2]) | uint16(b8[3])<<8
	seq := uint32(b8[4]) | uint32(b8[5])<<8 | uint32(b8[6])<<16 | uint32(b8[7])<<24
	return n, channelID, seq, nil
}

func (c *Client) handlePackageHeader(r io.Reader) error {
	_, channelID, seq, err := readPackageHeader(r)

	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	lastSeq, ok := c.getSeqNum(channelID)
	if !ok {
		c.setSeqNum(channelID, seq)
		return nil
	}

	// jumped pkg
	if (seq > lastSeq && seq-lastSeq > 1) || (lastSeq == math.MaxUint32 && seq > 0) {
		return fmt.Errorf("missing package, current seq: %d, last seq: %d", seq, lastSeq)
	}

	// connection reset
	if lastSeq < math.MaxInt32 && seq == 0 {
		return fmt.Errorf("connection reset, last seq: %d", lastSeq)
	}

	return nil
}

// Handle decodes an UDP packages into events.
// If it receives an InstrumentChange event, it has to update the instrumentIDToName accordingly.
// This function needs to handle for missing package, a jump in sequence number, e.g: prevSeq=5, curSeq=7
// And needs to handle for connection reset, the sequence number is zero.
// Note that the sequence number go from max(uint32) to zero is a normal event, not a connection reset.
func (c *Client) Handle(r io.Reader) error {
	err := c.handlePackageHeader(r)
	if err != nil {
		c.log.Errorw("failed to handle package header", "err", err)
		return err
	}

	events, err := c.decodeEvents(r)
	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	for _, event := range events {
		switch event.Type {
		case EventTypeInstrument:
			// update instrumentIDToName mapping
			var ins models.Instrument

			dataBytes, err := json.Marshal(event.Data)
			if err != nil {
				c.log.Errorw("failed to encode event data to []byte", "err", err)
				return err
			}

			err = json.Unmarshal(dataBytes, &ins)
			if err != nil {
				c.log.Errorw("failed to decode event data", "err", err)
				return err
			}
			c.mu.Lock()
			c.instrumentIDToName[ins.InstrumentID] = ins.InstrumentName
			c.mu.Unlock()
			c.Emit(newInstrumentNotificationChannel(), event)

		case EventTypeOrderBook:
			c.Emit(newOrderBookNotificationChannel(event.Instrument), event)

		case EventTypeTrades:
			c.Emit(newTradesNotificationChannel(event.Instrument), event)

		case EventTypeTicker:
			c.Emit(newTickerNotificationChannel(event.Instrument), event)

		default:
			return fmt.Errorf("invalid event type: %v", event.Type)
		}
	}

	return nil
}

func (c *Client) getConnection(addr string) *net.UDPConn {
	return c.addrConnMap[addr]
}

// ListenToEventsForAddress listens to one udp address on given network interface.
// If it receives an InstrumentChange event, it has to update the instrumentIDToName accordingly.
func (c *Client) ListEventsForAddress(addr string) error {
	conn := c.getConnection(addr)
	defer conn.Close()

	// listen to event and then handle
	for {
		data := make([]byte, defaultReadBufferSize)
		n, err := conn.Read(data)
		if err != nil {
			c.log.Errorw("fail to read UDP package", "error", err)
			return err
		} else {
			toBuffer := bytes.NewBuffer(data[:n])
			err := c.Handle(toBuffer)
			c.log.Errorw("fail to handle UDP package", "error", err)
			return err
		}

	}
}

// ListenToEvents listens to a list of udp addresses on given network interface.
// First it needs to build instrumentIDToName mapping using for mapping instrument id to instrument name.
func (c *Client) ListenToEvents() error {
	for _, addr := range c.addrs {
		go c.ListEventsForAddress(addr)
	}
	return nil
}
