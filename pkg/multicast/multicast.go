package multicast

import (
	"bytes"
	"context"
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
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
)

const (
	maxPacketSize     = 1500
	defaultDataChSize = 1000
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
	log      *zap.SugaredLogger
	inf      *net.Interface
	ipAddrs  []string
	port     int
	conn     *ipv4.PacketConn
	wsClient *websocket.Client

	supportCurrencies []string
	instrumentsMap    map[uint32]models.Instrument
	emitter           *emission.Emitter
}

// NewClient creates a new Client instance.
func NewClient(ifname string, ipAddrs []string, port int, wsClient *websocket.Client, currencies []string) (client *Client, err error) {
	l := zap.S()

	var inf *net.Interface
	if ifname != "" {
		inf, err = net.InterfaceByName(ifname)
		if err != nil {
			l.Errorw("failed to create net interfaces by name", "err", err, "ifname", ifname)
			return nil, err
		}
	}

	client = &Client{
		log:      l,
		inf:      inf,
		ipAddrs:  ipAddrs,
		port:     port,
		wsClient: wsClient,

		supportCurrencies: currencies,
		instrumentsMap:    make(map[uint32]models.Instrument),
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
		c.log.Errorw("failed to decode instrument event", "err", err)
		return Event{}, err
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
		c.log.Errorw("failed to decode orderbook event", "err", err)
		return Event{}, err
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
		c.log.Errorw("failed to decode trades event", "err", err)
		return Event{}, err
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
		c.log.Errorw("failed to decode ticker event", "err", err)
		return Event{}, err
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

func (c *Client) emitEvents(events []Event) {
	for _, event := range events {
		switch event.Type {
		case EventTypeInstrument:
			// update instrumentsMap
			ins := event.Data.(models.Instrument)
			c.setInstrument(ins.InstrumentID, ins)

			// emit event
			c.Emit(newInstrumentNotificationChannel(ins.Kind, ins.BaseCurrency), &ins)
			c.Emit(newInstrumentNotificationChannel(KindAny, ins.BaseCurrency), &ins)

		case EventTypeOrderBook:
			books := event.Data.(models.OrderBookRawNotification)
			c.Emit(newOrderBookNotificationChannel(books.InstrumentName), &books)

		case EventTypeTrades:
			trades := event.Data.(models.TradesNotification)
			if len(trades) > 0 {
				tradeIns := trades[0].InstrumentName
				tradeKind := trades[0].InstrumentKind
				currency := getCurrencyFromInstrument(tradeIns)
				c.Emit(newTradesNotificationChannel(tradeKind, currency), &trades)
			}

		case EventTypeTicker:
			ticker := event.Data.(models.TickerNotification)
			c.Emit(newTickerNotificationChannel(ticker.InstrumentName), &ticker)
		}
	}
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

// restartConnection should re-build instrumentMap.
func (c *Client) restartConnections(ctx context.Context) error {
	err := c.Stop()
	if err != nil {
		return err
	}
	return c.Start(ctx)
}

// Stop stops listening for events.
func (c *Client) Stop() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Client) getInstrument(id uint32) models.Instrument {
	return c.instrumentsMap[id]
}

func (c *Client) setInstrument(id uint32, ins models.Instrument) {
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

func (c *Client) handlePackageHeader(r io.Reader, chanelIDSeq map[uint16]uint32) error {
	_, channelID, seq, err := readPackageHeader(r)

	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	lastSeq, ok := chanelIDSeq[channelID]

	if ok {
		l := c.log.With("channelID", channelID, "current_seq", seq, "last_seq", lastSeq)
		if seq == 0 && math.MaxUint32-lastSeq >= 2 {
			return ErrConnectionReset
		}
		if seq == lastSeq {
			l.Debugw("package duplicated")
			return ErrDuplicatedPackage
		}
		if seq != lastSeq+1 {
			l.Warnw("package out of order")
			// return ErrOutOfOrder
		}
	}

	chanelIDSeq[channelID] = seq
	return nil
}

// Handle decodes an UDP packages into events.
// If it receives an InstrumentChange event, it has to update the instrumentsMap accordingly.
// This function needs to handle for missing package, a jump in sequence number, e.g: prevSeq=5, curSeq=7
// And needs to handle for connection reset, the sequence number is zero.
// Note that the sequence number go from max(uint32) to zero is a normal event, not a connection reset.
func (c *Client) Handle(m *sbe.SbeGoMarshaller, r io.Reader, chanelIDSeq map[uint16]uint32) error {
	err := c.handlePackageHeader(r, chanelIDSeq)
	if err != nil {
		if errors.Is(err, ErrDuplicatedPackage) {
			return nil
		}
		c.log.Errorw("failed to handle package header", "err", err)
		return err
	}

	events, err := c.decodeEvents(m, r)
	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	c.emitEvents(events)
	return nil
}

type Bytes []byte
type Pool struct {
	mu            *sync.RWMutex
	queue         []Bytes
	maxPacketSize int
}

func NewPool(maxPacketSize int) *Pool {
	q := make([]Bytes, 0)

	return &Pool{
		mu:            &sync.RWMutex{},
		queue:         q,
		maxPacketSize: maxPacketSize,
	}
}

func (p *Pool) Get() Bytes {
	p.mu.Lock()
	defer p.mu.Unlock()

	if n := len(p.queue); n > 0 {
		item := p.queue[n-1]
		p.queue[n-1] = nil
		p.queue = p.queue[:n-1]
		return item
	}
	return make(Bytes, p.maxPacketSize)
}

func (p *Pool) Put(bs Bytes) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.queue = append(p.queue, bs[:cap(bs)])
}

func (c *Client) setupConnection() ([]net.IP, error) {
	packetConn, err := net.ListenPacket("udp4", fmt.Sprintf("0.0.0.0:%d", c.port))

	if err != nil {
		c.log.Errorw("failed to initiate packet connection", "err", err, "port", c.port)
		return nil, err
	}

	c.conn = ipv4.NewPacketConn(packetConn)

	ipGroups := make([]net.IP, len(c.ipAddrs))

	err = c.conn.SetControlMessage(ipv4.FlagDst, true)
	if err != nil {
		c.log.Errorw("failed to set control message", "err", err)
		return nil, err
	}

	for index, ipAddr := range c.ipAddrs {
		group := net.ParseIP(ipAddr)
		if group == nil {
			return nil, ErrInvalidIpv4Address
		}
		err := c.conn.JoinGroup(c.inf, &net.UDPAddr{IP: group})
		if err != nil {
			c.log.Errorw("failed to join group", "group", group, "err", err, "ipAddr", ipAddr)
			return nil, err
		}
		ipGroups[index] = group
	}

	return ipGroups, nil
}

// ListenToEvents listens to a list of udp addresses on given network interface.
func (c *Client) ListenToEvents(ctx context.Context) error {
	ipGroups, err := c.setupConnection()
	if err != nil {
		c.log.Errorw("failed to setup ipv4 packet connection", "err", err)
		return nil
	}

	dataCh := make(chan []byte, defaultDataChSize)
	pool := NewPool(maxPacketSize)

	// handle data from dataCh
	go func() {
		m := sbe.NewSbeGoMarshaller()
		channelIDSeq := make(map[uint16]uint32)
		for {
			select {
			case <-ctx.Done():
				c.conn.Close()
			case data, ok := <-dataCh:
				if !ok {
					return
				}
				bufferData := bytes.NewBuffer(data)
				err := c.Handle(m, bufferData, channelIDSeq)
				if errors.Is(err, ErrConnectionReset) {
					c.log.Infow("restarting connection...", "error", err)
					err := c.restartConnections(ctx)
					if err != nil {
						c.log.Error("failed to restart connections", "err", err)
					}
				} else if err != nil {
					c.log.Errorw("fail to handle UDP package", "error", err)
				}
				pool.Put(data)
			}
		}
	}()

	// listen to event using ipv4 package
	go func() {
		for {
			data := pool.Get()
			n, cm, _, err := c.conn.ReadFrom(data)
			if err != nil {
				if isNetConnClosedErr(err) {
					c.log.Infow("connection closed", "error", err)
					close(dataCh)
					break
				}
			}

			if cm.Dst.IsMulticast() {
				if checkValidDstAddress(cm.Dst, ipGroups) { // joined group, push data to dataCh
					dataCh <- data[:n]
				} else { // unknown group, discard
					continue
				}
			}
		}
	}()

	return nil
}

func checkValidDstAddress(dest net.IP, groups []net.IP) bool {
	for _, ipGroup := range groups {
		if dest.Equal(ipGroup) {
			return true
		}
	}
	return false
}
