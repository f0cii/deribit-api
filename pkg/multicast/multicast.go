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

	RestartEventChannel = "multicast.restart"
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
func NewClient(
	ifname string,
	ipAddrs []string,
	port int,
	wsClient *websocket.Client,
	currencies []string,
) (client *Client, err error) {
	log := zap.S()

	var inf *net.Interface
	if ifname != "" {
		inf, err = net.InterfaceByName(ifname)
		if err != nil {
			log.Errorw("failed to create net interfaces by name", "err", err, "ifname", ifname)
			return nil, err
		}
	}

	client = &Client{
		log:      log,
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
func getAllInstrument(
	wsClient *websocket.Client, currencies []string,
) ([]models.Instrument, error) {
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
func (c *Client) decodeEvents(
	marshaler *sbe.SbeGoMarshaller, reader io.Reader,
) (events []Event, err error) {
	// Decode Sbe messages
	for {
		// Decode message header.
		var header sbe.MessageHeader
		err = header.Decode(marshaler, reader)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return
		}
		event, err := c.decodeEvent(marshaler, reader, header)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
}

func (c *Client) decodeEvent(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	header sbe.MessageHeader,
) (Event, error) {
	switch header.TemplateId {
	case 1000:
		return c.decodeInstrumentEvent(marshaler, reader, header)
	case 1001:
		return c.decodeOrderBookEvent(marshaler, reader, header)
	case 1002:
		return c.decodeTradesEvent(marshaler, reader, header)
	case 1003:
		return c.decodeTickerEvent(marshaler, reader, header)
	default:
		return Event{}, fmt.Errorf("%w, templateId: %d", ErrUnsupportedTemplateID, header.TemplateId)
	}
}

func (c *Client) decodeInstrumentEvent(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	header sbe.MessageHeader,
) (Event, error) {
	var ins sbe.Instrument
	err := ins.Decode(marshaler, reader, header.BlockLength, true)
	if err != nil {
		c.log.Errorw("failed to decode instrument event", "err", err)
		return Event{}, err
	}

	instrument := models.Instrument{
		TickSize:             ins.TickSize,
		TakerCommission:      ins.TakerCommission,
		SettlementPeriod:     ins.SettlementPeriod.String(),
		QuoteCurrency:        getCurrencyFromBytesArray(ins.QuoteCurrency),
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
		BaseCurrency:         getCurrencyFromBytesArray(ins.BaseCurrency),
		BlockTradeCommission: ins.BlockTradeCommission,
		OptionType:           ins.OptionType.String(),
		Strike:               ins.StrikePrice,
	}

	return Event{
		Type: EventTypeInstrument,
		Data: instrument,
	}, nil
}

func (c *Client) decodeOrderBookEvent(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	header sbe.MessageHeader,
) (Event, error) {
	var book sbe.Book
	err := book.Decode(marshaler, reader, header.BlockLength, true)
	if err != nil {
		c.log.Errorw("failed to decode orderbook event", "err", err)
		return Event{}, err
	}

	instrumentName := c.getInstrument(book.InstrumentId).InstrumentName
	event := models.OrderBookRawNotification{
		Timestamp:      int64(book.TimestampMs),
		InstrumentName: instrumentName,
		PrevChangeID:   int64(book.PrevChangeId),
		ChangeID:       int64(book.ChangeId),
	}

	for _, bookChange := range book.ChangesList {
		item := models.OrderBookNotificationItem{
			Action: bookChange.Change.String(),
			Price:  bookChange.Price,
			Amount: bookChange.Amount,
		}

		if bookChange.Side == sbe.BookSide.Ask {
			event.Asks = append(event.Asks, item)
		} else if bookChange.Side == sbe.BookSide.Bid {
			event.Bids = append(event.Bids, item)
		}
	}

	return Event{
		Type: EventTypeOrderBook,
		Data: event,
	}, nil
}

func (c *Client) decodeTradesEvent(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	header sbe.MessageHeader,
) (Event, error) {
	var trades sbe.Trades
	err := trades.Decode(marshaler, reader, header.BlockLength, true)
	if err != nil {
		c.log.Errorw("failed to decode trades event", "err", err)
		return Event{}, err
	}

	ins := c.getInstrument(trades.InstrumentId)

	tradesEvent := make(models.TradesNotification, len(trades.TradesList))
	for i, trade := range trades.TradesList {
		tradesEvent[i] = models.Trade{
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
		Data: tradesEvent,
	}, nil
}

func (c *Client) decodeTickerEvent(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	header sbe.MessageHeader,
) (Event, error) {
	var ticker sbe.Ticker
	err := ticker.Decode(marshaler, reader, header.BlockLength, true)
	if err != nil {
		c.log.Errorw("failed to decode ticker event", "err", err)
		return Event{}, err
	}

	instrumentName := c.getInstrument(ticker.InstrumentId).InstrumentName

	event := models.TickerNotification{
		Timestamp:       ticker.TimestampMs,
		Stats:           models.Stats{},
		State:           ticker.InstrumentState.String(),
		SettlementPrice: ticker.SettlementPrice,
		OpenInterest:    ticker.OpenInterest,
		MinPrice:        ticker.MinSellPrice,
		MaxPrice:        ticker.MaxBuyPrice,
		MarkPrice:       ticker.MarkPrice,
		LastPrice:       ticker.LastPrice,
		InstrumentName:  instrumentName,
		IndexPrice:      ticker.IndexPrice,
		Funding8H:       ticker.Funding8h,
		CurrentFunding:  ticker.CurrentFunding,
		BestBidPrice:    &ticker.BestBidPrice,
		BestBidAmount:   ticker.BestBidAmount,
		BestAskPrice:    &ticker.BestAskPrice,
		BestAskAmount:   ticker.BestAskAmount,
	}

	return Event{
		Type: EventTypeTicker,
		Data: event,
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

	err = c.Start(ctx)
	if err != nil {
		return err
	}

	c.Emit(RestartEventChannel, true)
	return nil
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

func readPackageHeader(reader io.Reader) (uint16, uint16, uint32, error) {
	b8 := make([]byte, 8)
	if _, err := io.ReadFull(reader, b8); err != nil {
		return 0, 0, 0, err
	}

	n := uint16(b8[0]) | uint16(b8[1])<<8
	channelID := uint16(b8[2]) | uint16(b8[3])<<8
	seq := uint32(b8[4]) | uint32(b8[5])<<8 | uint32(b8[6])<<16 | uint32(b8[7])<<24
	return n, channelID, seq, nil
}

func (c *Client) handlePackageHeader(reader io.Reader, chanelIDSeq map[uint16]uint32) error {
	_, channelID, seq, err := readPackageHeader(reader)
	if err != nil {
		c.log.Errorw("failed to decode events", "err", err)
		return err
	}

	lastSeq, ok := chanelIDSeq[channelID]

	if ok {
		log := c.log.With("channelID", channelID, "current_seq", seq, "last_seq", lastSeq)
		if seq == 0 && math.MaxUint32-lastSeq >= 2 {
			return ErrConnectionReset
		}
		if seq == lastSeq {
			log.Debugw("package duplicated")
			return ErrDuplicatedPackage
		}
		if seq != lastSeq+1 {
			log.Warnw("package out of order")
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
func (c *Client) Handle(
	marshaler *sbe.SbeGoMarshaller,
	reader io.Reader,
	chanelIDSeq map[uint16]uint32,
) error {
	err := c.handlePackageHeader(reader, chanelIDSeq)
	if err != nil {
		if errors.Is(err, ErrDuplicatedPackage) {
			return nil
		}
		c.log.Errorw("failed to handle package header", "err", err)
		return err
	}

	events, err := c.decodeEvents(marshaler, reader)
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
	lc := net.ListenConfig{
		Control: Control,
	}
	conn, err := lc.ListenPacket(context.Background(), "udp4", "0.0.0.0:"+strconv.Itoa(c.port))
	if err != nil {
		c.log.Errorw("Failed to initiate packet connection", "err", err, "port", c.port)
		return nil, err
	}

	c.conn = ipv4.NewPacketConn(conn)

	ipGroups := make([]net.IP, len(c.ipAddrs))

	err = c.conn.SetControlMessage(ipv4.FlagDst, true)
	if err != nil {
		c.log.Errorw("Failed to set control message", "err", err)
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
// nolint:cyclop
func (c *Client) ListenToEvents(ctx context.Context) error {
	ipGroups, err := c.setupConnection()
	if err != nil {
		c.log.Errorw("failed to setup ipv4 packet connection", "err", err)
		return err
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
				err := c.handleUDPPackage(ctx, m, channelIDSeq, data)
				if err != nil {
					c.log.Errorw("Fail to handle UDP package", "error", err)
				}
				pool.Put(data)
			}
		}
	}()

	// listen to event using ipv4 package
	go func() {
		for {
			data := pool.Get()

			res, err := readUDPMulticastPackage(c.conn, ipGroups, data)
			if res == nil {
				pool.Put(data)
			}

			if err != nil {
				if isNetConnClosedErr(err) {
					c.log.Infow("Connection closed", "error", err)
					close(dataCh)
					break
				}
				c.log.Errorw("Fail to read UDP multicast package", "error", err)
			} else if res != nil {
				dataCh <- res
			}
		}
	}()

	return nil
}

func (c *Client) handleUDPPackage(
	ctx context.Context,
	m *sbe.SbeGoMarshaller,
	channelIDSeq map[uint16]uint32,
	data []byte,
) error {
	buf := bytes.NewBuffer(data)
	err := c.Handle(m, buf, channelIDSeq)
	if err != nil {
		if errors.Is(err, ErrConnectionReset) {
			if err = c.restartConnections(ctx); err != nil {
				c.log.Error("failed to restart connections", "err", err)
				return err
			}
		} else {
			c.log.Errorw("Fail to handle UDP package", "error", err)
			return err
		}
	}

	return nil
}

func readUDPMulticastPackage(conn *ipv4.PacketConn, ipGroups []net.IP, data []byte) ([]byte, error) {
	n, cm, _, err := conn.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	if cm.Dst.IsMulticast() {
		if checkValidDstAddress(cm.Dst, ipGroups) {
			return data[:n], nil
		}
	}

	return nil, nil
}

func checkValidDstAddress(dest net.IP, groups []net.IP) bool {
	for _, ipGroup := range groups {
		if dest.Equal(ipGroup) {
			return true
		}
	}
	return false
}
