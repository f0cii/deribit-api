package multicast

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/KyberNetwork/deribit-api/pkg/multicast/sbe"
	"github.com/KyberNetwork/deribit-api/pkg/websocket"
	"github.com/chuckpreslar/emission"
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
	Type EventType
	Data interface{}
}

// Client represents a client for Deribit multicast.
type Client struct {
	log         *zap.SugaredLogger
	mu          sync.RWMutex
	inf         *net.Interface
	addrs       []string
	addrConnMap map[string]*net.UDPConn
	wsClient    *websocket.Client

	supportCurrencies []string
	instrumentsMap    map[uint32]models.Instrument
	channelIDSeqNum   map[uint16]uint32
	emitter           *emission.Emitter
}

// NewClient creates a new Client instance.
func NewClient(ifname string, addrs []string, wsClient *websocket.Client, currencies []string) (client *Client, err error) {
	l := zap.S()

	var inf *net.Interface
	if ifname != "" {
		inf, err = net.InterfaceByName(ifname)
		if err != nil {
			l.Errorw("failed to create net interfaces by name", "err", err)
			return nil, err
		}
	}

	client = &Client{
		log:         l,
		mu:          sync.RWMutex{},
		inf:         inf,
		addrs:       addrs,
		wsClient:    wsClient,
		addrConnMap: make(map[string]*net.UDPConn),

		supportCurrencies: currencies,
		instrumentsMap:    make(map[uint32]models.Instrument),
		channelIDSeqNum:   make(map[uint16]uint32),
		emitter:           emission.NewEmitter(),
	}

	return client, nil
}

// buildInstrumentsMapping builds a mapping to map instrument id to instrument.
func (c *Client) buildInstrumentsMapping() error {
	// need to update this field via instruments_update event
	instruments, err := getAllInstrument(c.wsClient, c.supportCurrencies)
	if err != nil {
		c.log.Errorw("failed to get all instruments", "err", err)
		return err
	}
	for _, instrument := range instruments {
		c.instrumentsMap[instrument.InstrumentID] = instrument
	}
	return nil
}

// getAllInstrument returns a list of all instruments by currencies
func getAllInstrument(wsClient *websocket.Client, currencies []string) ([]models.Instrument, error) {
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

// decodeEvents decodes a UDP package into a list of events.
func (c *Client) decodeEvents(m *sbe.SbeGoMarshaller, r io.Reader) (events []Event, err error) {
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
		event, err := c.decodeEvent(m, r, header)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
}

func (c *Client) decodeEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	switch header.TemplateId {
	case 1000:
		return c.decodeInstrumentEvent(m, r, header)
	case 1001:
		return c.decodeOrderBookEvent(m, r, header)
	case 1002:
		return c.decodeTradesEvent(m, r, header)
	case 1003:
		return c.decodeTickerEvent(m, r, header)
	default:
		return Event{}, fmt.Errorf("%w, templateId: %d", ErrUnsupportedTemplateId, header.TemplateId)
	}
}

func (c *Client) decodeInstrumentEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var ins sbe.Instrument
	err := ins.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	instrument := models.Instrument{
		TickSize:             ins.TickSize,
		TakerCommission:      ins.TakerCommission,
		SettlementPeriod:     ins.SettlementPeriod.String(),
		QuoteCurrency:        string(ins.QuoteCurrency[:]),
		MinTradeAmount:       ins.MinTradeAmount,
		MakerCommission:      ins.MakerCommission,
		Leverage:             int(ins.MaxLeverage),
		Kind:                 ins.Kind.String(),
		IsActive:             ins.InstrumentState.IsActive(),
		InstrumentID:         ins.InstrumentId,
		InstrumentName:       string(ins.InstrumentName),
		ExpirationTimestamp:  ins.ExpirationTimestampMs,
		CreationTimestamp:    ins.CreationTimestampMs,
		ContractSize:         ins.ContractSize,
		BaseCurrency:         string(ins.BaseCurrency[:]),
		BlockTradeCommission: ins.BlockTradeCommission,
		OptionType:           ins.OptionType.String(),
		Strike:               ins.StrikePrice,
	}
	return Event{
		Type: EventTypeInstrument,
		Data: instrument,
	}, nil
}

func (c *Client) decodeOrderBookEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var b sbe.Book
	err := b.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	instrumentName := c.getInstrument(b.InstrumentId).InstrumentName
	book := models.OrderBookRawNotification{
		Timestamp:      int64(b.TimestampMs),
		InstrumentName: instrumentName,
		PrevChangeID:   int64(b.PrevChangeId),
		ChangeID:       int64(b.ChangeId),
	}

	for _, bookChange := range b.ChangesList {
		item := models.OrderBookNotificationItem{
			Action: bookChange.Change.String(),
			Price:  bookChange.Price,
			Amount: bookChange.Amount,
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

func (c *Client) decodeTradesEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var t sbe.Trades
	err := t.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	ins := c.getInstrument(t.InstrumentId)

	trades := make(models.TradesNotification, len(t.TradesList))
	for id, trade := range t.TradesList {

		trades[id] = models.Trade{
			Amount:         trade.Amount,
			BlockTradeID:   strconv.FormatUint(trade.BlockTradeId, 10),
			Direction:      trade.Direction.String(),
			IndexPrice:     trade.IndexPrice,
			InstrumentName: ins.InstrumentName,
			InstrumentKind: ins.Kind,
			IV:             trade.Iv,
			Liquidation:    trade.Liquidation.String(),
			MarkPrice:      trade.MarkPrice,
			Price:          trade.Price,
			TickDirection:  int(trade.TickDirection),
			Timestamp:      trade.TimestampMs,
			TradeID:        strconv.FormatUint(trade.BlockTradeId, 10),
			TradeSeq:       trade.TradeSeq,
		}
	}

	return Event{
		Type: EventTypeTrades,
		Data: trades,
	}, nil
}

func (c *Client) decodeTickerEvent(m *sbe.SbeGoMarshaller, r io.Reader, header sbe.MessageHeader) (Event, error) {
	var t sbe.Ticker
	err := t.Decode(m, r, header.BlockLength, true)
	if err != nil {
		return Event{}, nil
	}

	instrumentName := c.getInstrument(t.InstrumentId).InstrumentName

	ticker := models.TickerNotification{
		Timestamp:       t.TimestampMs,
		Stats:           models.Stats{},
		State:           t.InstrumentState.String(),
		SettlementPrice: t.SettlementPrice,
		OpenInterest:    t.OpenInterest,
		MinPrice:        t.MinSellPrice,
		MaxPrice:        t.MaxBuyPrice,
		MarkPrice:       t.MarkPrice,
		LastPrice:       t.LastPrice,
		InstrumentName:  instrumentName,
		IndexPrice:      t.IndexPrice,
		Funding8H:       t.Funding8h,
		CurrentFunding:  t.CurrentFunding,
		BestBidPrice:    &t.BestBidPrice,
		BestBidAmount:   t.BestBidAmount,
		BestAskPrice:    &t.BestAskPrice,
		BestAskAmount:   t.BestAskAmount,
	}
	return Event{
		Type: EventTypeTicker,
		Data: ticker,
	}, nil
}

// Start starts listening events on interface `ifname` and `addrs`.
func (c *Client) Start(ctx context.Context) error {

	err := c.buildInstrumentsMapping()
	if err != nil {
		c.log.Errorw("failed to build instruments mapping", "err", err)
		return err
	}

	return c.ListenToEvents(ctx)
}

func (c *Client) closeConnection(addr string) error {
	conn := c.getConnection(addr)
	err := conn.Close()
	if err != nil {
		c.log.Errorw("failed to close connection", "err", err)
		return err
	}
	return nil
}

// restartConnection should re-build instrumentMap.
func (c *Client) restartConnection(ctx context.Context, addr string) error {
	conn := c.getConnection(addr)
	err := conn.Close()
	if err != nil {
		c.log.Errorw("failed to close connection", "err", err)
		return err
	}

	err = c.buildInstrumentsMapping()
	if err != nil {
		c.log.Errorw("failed to build instruments mapping", "err", err)
		return err
	}

	go c.ListenToEventsForAddress(ctx, addr)
	return nil
}

// Stop stops listening for events.
func (c *Client) Stop() error {
	for _, addr := range c.addrs {
		go c.closeConnection(addr)
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

func (c *Client) getInstrument(id uint32) models.Instrument {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.instrumentsMap[id]
}

func (c *Client) setInstrument(id uint32, ins models.Instrument) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.instrumentsMap[id] = ins
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

	// check for invalid sequence number
	if seq != lastSeq+1 {
		if seq == 0 {
			return ErrConnectionReset
		}
		return ErrLostPackage
	}

	return nil
}

// Handle decodes an UDP packages into events.
// If it receives an InstrumentChange event, it has to update the instrumentsMap accordingly.
// This function needs to handle for missing package, a jump in sequence number, e.g: prevSeq=5, curSeq=7
// And needs to handle for connection reset, the sequence number is zero.
// Note that the sequence number go from max(uint32) to zero is a normal event, not a connection reset.
func (c *Client) Handle(m *sbe.SbeGoMarshaller, r io.Reader) error {
	err := c.handlePackageHeader(r)
	if err != nil {
		c.log.Errorw("failed to handle package header", "err", err)
		return err
	}

	events, err := c.decodeEvents(m, r)
	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	for _, event := range events {
		switch event.Type {
		case EventTypeInstrument:
			// update instrumentsMap
			ins, ok := event.Data.(models.Instrument)
			if !ok {
				return fmt.Errorf("invalid event type")
			}
			c.setInstrument(ins.InstrumentID, ins)

			// emit event
			c.Emit(newInstrumentNotificationChannel(ins.Kind, ins.BaseCurrency), &event.Data)
			c.Emit(newInstrumentNotificationChannel(KindAny, ins.BaseCurrency), &event.Data)

		case EventTypeOrderBook:
			books := event.Data.(models.OrderBookRawNotification)
			c.Emit(newOrderBookNotificationChannel(books.InstrumentName), &event.Data)

		case EventTypeTrades:
			trades := event.Data.(models.TradesNotification)
			if len(trades) > 0 {
				tradeIns := trades[0].InstrumentName
				tradeKind := trades[0].InstrumentKind
				currency := getCurrencyFromInstrument(tradeIns)
				c.Emit(newTradesNotificationChannel(tradeKind, currency), &event.Data)
			}

		case EventTypeTicker:
			ticker := event.Data.(models.TickerNotification)
			c.Emit(newTickerNotificationChannel(ticker.InstrumentName), &event.Data)
		}
	}

	return nil
}

func (c *Client) getConnection(addr string) *net.UDPConn {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.addrConnMap[addr]
}

func (c *Client) setConnection(addr string, conn *net.UDPConn) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.addrConnMap[addr] = conn
}

func (c *Client) listenToMulticastUDP(addr string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		c.log.Errorw("failed to resolve UDP address", "err", err)
		return nil, err
	}

	udpConn, err := net.ListenMulticastUDP("udp", c.inf, udpAddr)
	if err != nil {
		c.log.Errorw("failed to listen to multicast UDP", "error", err)
		return nil, err
	}

	err = udpConn.SetReadBuffer(defaultReadBufferSize)
	if err != nil {
		c.log.Errorw("fail to set read buffer for UDP", "error", err)
		return nil, err
	}

	c.setConnection(addr, udpConn)
	return udpConn, nil
}

// ListenToEventsForAddress listens to one udp address on given network interface.
func (c *Client) ListenToEventsForAddress(ctx context.Context, addr string) error {
	dataCh := make(chan []byte, defaultDataChSize)
	udpConn, err := c.listenToMulticastUDP(addr)
	if err != nil {
		c.log.Errorw("failed to listen to multicast UDP", "err", err)
		return err
	}

	// handle data from dataCh
	go func() {
		m := sbe.NewSbeGoMarshaller()
		for {
			select {
			case <-ctx.Done():
				udpConn.Close()
			case data := <-dataCh:
				bufferData := bytes.NewBuffer(data)
				err := c.Handle(m, bufferData)
				if errors.Is(err, ErrConnectionReset) || errors.Is(err, ErrLostPackage) {
					c.log.Infow("connection reset or lost package err, restarting connection...", "error", err)
					err := c.restartConnection(ctx, addr)
					if err != nil {
						c.log.Errorw("fail to restart connection", "error", err)
					}
					return
				} else {
					c.log.Errorw("fail to handle UDP package", "error", err)
				}
			}
		}
	}()

	// listen to event
	for {
		data := make([]byte, defaultReadBufferSize)
		n, err := udpConn.Read(data)
		if err != nil {
			if isNetConnClosedErr(err) {
				c.log.Infow("connection closed", "error", err)
				break
			}
			c.log.Errorw("fail to read UDP package", "error", err)
		} else {
			dataCh <- data[:n]
		}

	}
	return nil
}

// ListenToEvents listens to a list of udp addresses on given network interface.
func (c *Client) ListenToEvents(ctx context.Context) error {
	for _, addr := range c.addrs {
		go c.ListenToEventsForAddress(ctx, addr)
	}
	return nil
}
