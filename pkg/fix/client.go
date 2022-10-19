package fix

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/chuckpreslar/emission"
	"github.com/google/uuid"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"go.uber.org/zap"
)

const (
	nonceLen = 64

	subscriptionChannelParts = 2
	subscriptionTypeBook     = "book"
	subscriptionTypeTrades   = "trades"
)

type Initiator interface {
	Start() error
	Stop()
}

type Sender func(m quickfix.Messagable) (err error)

type Config struct {
	APIKey    string
	SecretKey string
	Settings  *quickfix.Settings
	Dialer    Dialer
	Sender    Sender
}

// Client implements the quickfix.Application interface.
type Client struct {
	log *zap.SugaredLogger

	apiKey    string
	secretKey string

	settings *quickfix.Settings

	targetCompID string
	senderCompID string

	initiator Initiator

	mu          sync.Mutex
	isConnected bool

	sending sync.Mutex
	pending map[string]*call

	subscriptions    []string
	subscriptionsMap map[string]bool
	emitter          *emission.Emitter
	sender           Sender
}

type Dialer func(
	app quickfix.Application,
	storeFactory quickfix.MessageStoreFactory,
	appSettings *quickfix.Settings,
	logFactory quickfix.LogFactory,
) (Initiator, error)

// OnCreate implemented as part of Application interface.
func (c *Client) OnCreate(_ quickfix.SessionID) {}

// OnLogon implemented as part of Application interface.
func (c *Client) OnLogon(_ quickfix.SessionID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = true
	c.log.Debugw("Logon successfully!")
}

// OnLogout implemented as part of Application interface.
func (c *Client) OnLogout(_ quickfix.SessionID) {
	defer func() {
		if err := recover(); err != nil {
			c.log.Errorw("Recover from panic", "error", err)
		}
	}()

	c.mu.Lock()
	c.isConnected = false
	c.mu.Unlock()

	c.log.Debugw("Logged out!")
	for _, call := range c.pending {
		call.done <- ErrClosed
		close(call.done)
	}
}

// FromAdmin implemented as part of Application interface.
func (c *Client) FromAdmin(_ *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	return nil
}

// ToAdmin implemented as part of Application interface.
func (c *Client) ToAdmin(msg *quickfix.Message, _ quickfix.SessionID) {
	timestamp := time.Now().UnixMilli()
	nonce, err := generateRandomBytes(nonceLen)
	if err != nil {
		c.log.Errorw(
			"Fail to generate random bytes",
			"nonce_len", nonceLen,
			"error", err,
		)
		return
	}

	rawData := strconv.FormatInt(timestamp, 10) + "." + base64.StdEncoding.EncodeToString(nonce)
	hash := sha256.Sum256([]byte(rawData + c.secretKey))
	password := base64.StdEncoding.EncodeToString(hash[:])

	msg.Body.Set(field.NewRawData(rawData))
	msg.Body.Set(field.NewUsername(c.apiKey))
	msg.Body.Set(field.NewPassword(password))
	msg.Body.Set(field.NewResetSeqNumFlag(true))
	msg.Body.SetBool(tagCancelOnDisconnect, false)
}

// ToApp implemented as a part of Application interface.
func (c *Client) ToApp(msg *quickfix.Message, _ quickfix.SessionID) error {
	c.log.Debugw("Sending message to server", "msg", msg)
	return nil
}

// FromApp implemented as a part of Application interface.
func (c *Client) FromApp(msg *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	// Process message according to message type.
	msgType, err := msg.MsgType()
	if err != nil {
		c.log.Errorw("Fail to get response message type", "error", err)
		return err
	}

	c.handleSubscriptions(msgType, msg)

	reqIDTag, err2 := getReqIDTagFromMsgType(enum.MsgType(msgType))
	if err2 != nil {
		c.log.Warnw("Could not get request ID tag", "msgType", msgType, "error", err2)
		return nil
	}

	id, err := msg.Body.GetString(reqIDTag)
	if err != nil {
		c.log.Errorw("Fail to get request ID", "tag", reqIDTag, "error", err)
		return err
	}

	c.mu.Lock()
	call := c.pending[id]
	delete(c.pending, id)
	c.mu.Unlock()

	if call != nil {
		c.log.Debugw(
			"Matching response message",
			"id_tag", reqIDTag,
			"id", id,
			"request", call.request,
			"response", msg,
		)
		response, err2 := copyMessage(msg)
		if err2 != nil {
			c.log.Fatalw("Fail to copy response message", "error", err2)
		}
		call.response = response
		call.done <- nil
		close(call.done)
	}

	return nil
}

// New returns a new client for Deribit FIX API.
// nolint:funlen
func New(
	ctx context.Context,
	cfg Config,
) (*Client, error) {
	logger := zap.S()

	// Get TargetCompID and SenderCompID from settings.
	if cfg.Settings == nil {
		return nil, errors.New("empty quickfix settings")
	}

	globalSettings := cfg.Settings.GlobalSettings()
	targetCompID, err := globalSettings.Setting("TargetCompID")
	if err != nil {
		logger.Errorw("Fail to read TargetCompID from settings", "error", err)
		return nil, err
	}

	senderCompID, err := globalSettings.Setting("SenderCompID")
	if err != nil {
		logger.Errorw("Fail to read SenderCompID from settings", "error", err)
		return nil, err
	}
	sender := cfg.Sender
	if sender == nil {
		sender = quickfix.Send
	}

	// Create a new Client object.
	client := &Client{
		log:              logger,
		apiKey:           cfg.APIKey,
		secretKey:        cfg.SecretKey,
		settings:         cfg.Settings,
		targetCompID:     targetCompID,
		senderCompID:     senderCompID,
		mu:               sync.Mutex{},
		isConnected:      false,
		sending:          sync.Mutex{},
		pending:          make(map[string]*call),
		subscriptionsMap: make(map[string]bool),
		emitter:          emission.NewEmitter(),
		sender:           sender,
	}

	// Init session and logon to deribit FIX API server.
	logFactory := quickfix.NewNullLogFactory()

	dialer := cfg.Dialer
	if dialer == nil {
		dialer = func(
			a quickfix.Application,
			f quickfix.MessageStoreFactory,
			s *quickfix.Settings,
			l quickfix.LogFactory,
		) (Initiator, error) {
			return quickfix.NewInitiator(a, f, s, l)
		}
	}

	client.initiator, err = dialer(
		client,
		quickfix.NewMemoryStoreFactory(),
		cfg.Settings,
		logFactory,
	)
	if err != nil {
		client.log.Errorw("Fail to create new initiator", "error", err)
		return nil, err
	}

	err = client.Start()
	if err != nil {
		client.log.Errorw("Fail to start fix connection", "error", err)
		return nil, err
	}

	return client, nil
}

func (c *Client) Start() error {
	c.mu.Lock()
	c.subscriptionsMap = make(map[string]bool)
	c.mu.Unlock()

	if err := c.initiator.Start(); err != nil {
		c.log.Errorw("Fail to initialize initiator", "error", err)
		return err
	}

	// Wait for the session to be authorized by the server.
	for !c.IsConnected() {
		time.Sleep(10 * time.Millisecond)
	}

	if len(c.subscriptions) > 0 {
		err := c.Subscribe(context.Background(), c.subscriptions)
		if err != nil {
			c.log.Warnw("Fail to resubscribe to channels", "error", err)
		}
	}

	return nil
}

// IsConnected checks whether the connection is established or not.
func (c *Client) IsConnected() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.isConnected
}

// Close closes underlying connection.
func (c *Client) Close() {
	c.initiator.Stop()
}

// nolint:funlen,gocognit,cyclop
func (c *Client) handleSubscriptions(msgType string, msg *quickfix.Message) {
	logger := c.log.With("msg", msg)

	switch enum.MsgType(msgType) {
	case enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH, enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH:
		symbol, err := getSymbol(msg)
		if err != nil {
			logger.Warnw("Fail to get symbol", "error", err)
			return
		}

		entries, err := getMDEntries(msg)
		if err != nil {
			logger.Warnw("Fail to get NoMDEntries", "error", err)
			return
		}

		markPrice, err := getMarkPrice(msg)
		if err != nil {
			logger.Warnw("Fail to get mark price", "error", err)
			return
		}

		var tradesEvent models.TradesNotification
		orderBookEvent := models.OrderBookRawNotification{
			InstrumentName: symbol,
		}
		for i := 0; i < entries.Len(); i++ {
			entry := entries.Get(i)
			entryType, err := getMDEntryType(entry)
			if err != nil {
				logger.Warnw("No value for MDEntryType", "error", err)
				continue
			}

			if entryType != enum.MDEntryType_BID &&
				entryType != enum.MDEntryType_OFFER &&
				entryType != enum.MDEntryType_TRADE {
				continue
			}

			serverTime, err := getMDEntryDate(entry)
			if err != nil {
				logger.Warnw("Fail to get MDEntryTime", "error", err)
				continue
			}

			price, err := getMDEntryPx(entry)
			if err != nil {
				logger.Warnw("Fail to get MDEntryPx", "error", err)
				continue
			}

			var action string
			if !hasMDUpdateAction(entry) {
				action = "new"
			} else {
				action, err = getMDUpdateAction(entry)
				if err != nil {
					logger.Warnw("Fail to get MDUpdateAction", "error", err)
					continue
				}
			}

			amount, err := getMDEntrySize(entry)
			if err != nil {
				logger.Warnw("Fail to get MDEntrySize", "error", err)
				continue
			}

			switch entryType {
			case enum.MDEntryType_BID:
				item := models.OrderBookNotificationItem{
					Action: action,
					Price:  price,
					Amount: amount,
				}
				orderBookEvent.Bids = append(orderBookEvent.Bids, item)
				orderBookEvent.Timestamp = serverTime.UnixMilli()
			case enum.MDEntryType_OFFER:
				item := models.OrderBookNotificationItem{
					Action: action,
					Price:  price,
					Amount: amount,
				}
				orderBookEvent.Asks = append(orderBookEvent.Asks, item)
				orderBookEvent.Timestamp = serverTime.UnixMilli()
			case enum.MDEntryType_TRADE:
				indexPrice, err := getGroupPrice(entry)
				if err != nil {
					logger.Warnw("Fail to get index price", "error", err)
					continue
				}

				side, err := getGroupSide(entry)
				if err != nil {
					logger.Warnw("Fail to get trade side", "error", err)
					continue
				}

				tradeID, err := getGroupTradeID(entry)
				if err != nil {
					logger.Warnw("Fail to get trade ID", "error", err)
					continue
				}

				trade := models.Trade{
					Amount:         amount,
					Direction:      decodeOrderSide(side),
					IndexPrice:     indexPrice,
					InstrumentName: symbol,
					MarkPrice:      markPrice,
					Price:          price,
					Timestamp:      uint64(serverTime.UnixMilli()),
					TradeID:        tradeID,
				}
				tradesEvent = append(tradesEvent, trade)
			}
		}

		if len(orderBookEvent.Bids) > 0 || len(orderBookEvent.Asks) > 0 {
			var isSnapshot bool
			if msgType == string(enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH) {
				isSnapshot = true
			}

			c.Emit(newOrderBookNotificationChannel(symbol), &orderBookEvent, isSnapshot)
		}

		if len(tradesEvent) > 0 {
			c.Emit(newTradeNotificationChannel(symbol), &tradesEvent)
		}
	default:
		return
	}
}

func (c *Client) addCommonHeaders(msg *quickfix.Message) {
	msg.Header.Set(field.NewBeginString(fixVersion))
	msg.Header.Set(field.NewTargetCompID(c.targetCompID))
	msg.Header.Set(field.NewSenderCompID(c.senderCompID))
	msg.Header.Set(field.NewSendingTime(time.Now().UTC()))
}

func (c *Client) send(
	_ context.Context, id string, msg *quickfix.Message, wait bool,
) (Waiter, error) {
	c.sending.Lock()
	defer c.sending.Unlock()

	c.mu.Lock()
	if !c.isConnected {
		c.mu.Unlock()
		return Waiter{}, ErrClosed
	}

	c.addCommonHeaders(msg)
	var cc *call
	if wait {
		cc = &call{request: msg, done: make(chan error, 1)}
		c.pending[id] = cc
	}
	c.mu.Unlock()

	if err := c.sender(msg); err != nil {
		c.mu.Lock()
		delete(c.pending, id)
		c.mu.Unlock()
		return Waiter{}, err
	}

	return Waiter{call: cc}, nil
}

// Call initiates a FIX call and wait for the response.
func (c *Client) Call(
	ctx context.Context, id string, msg *quickfix.Message,
) (*quickfix.Message, error) {
	call, err := c.send(ctx, id, msg, true)
	if err != nil {
		return nil, err
	}

	return call.Wait(ctx)
}

type call struct {
	request  *quickfix.Message
	response *quickfix.Message
	done     chan error
}

// Waiter proxies an ongoing FIX call.
type Waiter struct {
	*call
}

// Wait for the response message of an ongoing FIX call.
func (w Waiter) Wait(ctx context.Context) (*quickfix.Message, error) {
	select {
	case err, ok := <-w.call.done:
		if !ok {
			err = ErrClosed
		}
		if err != nil {
			return nil, err
		}
		return w.call.response, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Client) MarketDataRequest(
	ctx context.Context,
	subscriptionRequestType enum.SubscriptionRequestType,
	marketDepth *int,
	mdUpdateType *enum.MDUpdateType,
	mdEntryTypes []enum.MDEntryType,
	instruments []string,
) (*quickfix.Message, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		c.log.Errorw("Fail to generate uuid", "error", err)
		return nil, err
	}

	msg := quickfix.NewMessage()
	msg.Header.Set(field.NewMsgType(enum.MsgType_MARKET_DATA_REQUEST))

	msg.Body.Set(field.NewMDReqID(id.String()))
	msg.Body.Set(field.NewSubscriptionRequestType(subscriptionRequestType))
	if marketDepth != nil {
		msg.Body.Set(field.NewMarketDepth(*marketDepth))
	}
	if mdUpdateType != nil {
		msg.Body.Set(field.NewMDUpdateType(*mdUpdateType))
	}

	if len(mdEntryTypes) > 0 {
		groups := newNoMDEntryTypesRepeatingGroup()
		for _, entryType := range mdEntryTypes {
			g := groups.Add()
			g.Set(field.NewMDEntryType(entryType))
		}
		msg.Body.SetGroup(groups)
	}

	if len(instruments) > 0 {
		groups := newNoRelatedSymRepeatingGroup()
		for _, instrument := range instruments {
			g := groups.Add()
			g.Set(field.NewSymbol(instrument))
		}
		msg.Body.SetGroup(groups)
	}

	return c.Call(ctx, id.String(), msg)
}

func (c *Client) SubscribeOrderBooks(ctx context.Context, instruments []string) error {
	if len(instruments) == 0 {
		c.log.Debugw("No instruments to subscribe")
		return nil
	}

	marketDepth := 0
	mdUpdateType := enum.MDUpdateType_INCREMENTAL_REFRESH
	msg, err := c.MarketDataRequest(
		ctx,
		enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
		&marketDepth,
		&mdUpdateType,
		[]enum.MDEntryType{
			enum.MDEntryType_BID,
			enum.MDEntryType_OFFER,
		},
		instruments,
	)
	if err != nil {
		c.log.Errorw("Fail to subscribe orderbooks", "error", err)
		return err
	}

	if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_REQUEST_REJECT)) {
		reason, err := getText(msg)
		if err != nil {
			c.log.Warnw("No value for Text field", "error", err)
		} else {
			err = errors.New(reason)
		}
		return err
	}

	return nil
}

func (c *Client) UnsubscribeOrderBooks(ctx context.Context, instruments []string) error {
	if len(instruments) == 0 {
		c.log.Debugw("No instruments to unsubscribe")
		return nil
	}

	msg, err := c.MarketDataRequest(
		ctx,
		enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST,
		nil,
		nil,
		[]enum.MDEntryType{
			enum.MDEntryType_BID,
			enum.MDEntryType_OFFER,
		},
		instruments,
	)
	if err != nil {
		c.log.Errorw("Fail to unsubscribe orderbooks", "error", err)
		return err
	}

	if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_REQUEST_REJECT)) {
		reason, err := getText(msg)
		if err != nil {
			c.log.Warnw("No value for Text field", "error", err)
		} else {
			err = errors.New(reason)
		}
		return err
	}

	return nil
}

func (c *Client) SubscribeTrades(ctx context.Context, instruments []string) error {
	if len(instruments) == 0 {
		c.log.Debugw("No instruments to subscribe")
		return nil
	}

	marketDepth := 1
	mdUpdateType := enum.MDUpdateType_INCREMENTAL_REFRESH
	msg, err := c.MarketDataRequest(
		ctx,
		enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
		&marketDepth,
		&mdUpdateType,
		[]enum.MDEntryType{
			enum.MDEntryType_TRADE,
		},
		instruments,
	)
	if err != nil {
		c.log.Errorw("Fail to subscribe trades", "error", err)
		return err
	}

	if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_REQUEST_REJECT)) {
		reason, err := getText(msg)
		if err != nil {
			c.log.Warnw("No value for Text field", "error", err)
		} else {
			err = errors.New(reason)
		}
		return err
	}

	return nil
}

func (c *Client) UnsubscribeTrades(ctx context.Context, instruments []string) error {
	if len(instruments) == 0 {
		c.log.Debugw("No instruments to unsubscribe")
		return nil
	}

	msg, err := c.MarketDataRequest(
		ctx,
		enum.SubscriptionRequestType_DISABLE_PREVIOUS_SNAPSHOT_PLUS_UPDATE_REQUEST,
		nil,
		nil,
		[]enum.MDEntryType{
			enum.MDEntryType_TRADE,
		},
		instruments,
	)
	if err != nil {
		c.log.Errorw("Fail to unsubscribe trades", "error", err)
		return err
	}

	if msg.IsMsgTypeOf(string(enum.MsgType_MARKET_DATA_REQUEST_REJECT)) {
		reason, err := getText(msg)
		if err != nil {
			c.log.Warnw("No value for Text field", "error", err)
		} else {
			err = errors.New(reason)
		}
		return err
	}

	return nil
}

// Subscribe listens for notifications.
// Currently, only support for orderbook notifications.
// nolint: cyclop
func (c *Client) Subscribe(ctx context.Context, channels []string) error {
	c.mu.Lock()
	instrumentMap := make(map[string][]string)
	for _, channel := range channels {
		if !c.subscriptionsMap[channel] {
			parts := strings.Split(channel, ".")
			if len(parts) != subscriptionChannelParts {
				continue // Ignore channels don't have format <SubType>.<Instrument>
			}

			if parts[0] != subscriptionTypeBook && parts[0] != subscriptionTypeTrades {
				continue // Support only book.<Instrument> and trades.<Instrument>
			}

			c.subscriptionsMap[channel] = true
			c.subscriptions = append(c.subscriptions, channel)
			if _, ok := instrumentMap[parts[0]]; !ok {
				instrumentMap[parts[0]] = []string{parts[1]}
			} else {
				instrumentMap[parts[0]] = append(instrumentMap[parts[0]], parts[1])
			}
		}
	}
	c.mu.Unlock()

	for subType, instruments := range instrumentMap {
		switch subType {
		case subscriptionTypeBook:
			err := c.SubscribeOrderBooks(ctx, instruments)
			if err != nil {
				c.log.Errorw("Fail to subscribe orderbook notifications", "error", err)
				return err
			}
		case subscriptionTypeTrades:
			err := c.SubscribeTrades(ctx, instruments)
			if err != nil {
				c.log.Errorw("Fail to subscribe trades notifications", "error", err)
				return err
			}
		}
	}

	return nil
}

// nolint: cyclop
func (c *Client) Unsubscribe(ctx context.Context, channels []string) error {
	c.mu.Lock()
	instrumentMap := make(map[string][]string)
	for _, channel := range channels {
		if c.subscriptionsMap[channel] {
			parts := strings.Split(channel, ".")
			if _, ok := instrumentMap[parts[0]]; !ok {
				instrumentMap[parts[0]] = []string{parts[1]}
			} else {
				instrumentMap[parts[0]] = append(instrumentMap[parts[0]], parts[1])
			}
		}
	}
	c.mu.Unlock()

	for subType, instruments := range instrumentMap {
		switch subType {
		case "book":
			err := c.UnsubscribeOrderBooks(ctx, instruments)
			if err != nil {
				c.log.Errorw("Fail to unsubscribe orderbook notifications", "error", err)
				return err
			}
		case "trades":
			err := c.UnsubscribeTrades(ctx, instruments)
			if err != nil {
				c.log.Errorw("Fail to unsubscribe trades notifications", "error", err)
				return err
			}
		}
	}

	c.mu.Lock()
	for _, channel := range channels {
		if c.subscriptionsMap[channel] {
			delete(c.subscriptionsMap, channel)
		}
	}
	c.subscriptions = c.subscriptions[:0]
	for channel := range c.subscriptionsMap {
		c.subscriptions = append(c.subscriptions, channel)
	}
	c.mu.Unlock()

	return nil
}

// nolint:funlen
func (c *Client) CreateOrder(
	ctx context.Context,
	instrument string,
	side enum.Side,
	amount float64,
	price float64,
	orderType enum.OrdType,
	timeInForce enum.TimeInForce,
	execInst string,
	label string,
) (order models.Order, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		c.log.Errorw("Fail to generate uuid", "error", err)
		return order, err
	}

	msg := quickfix.NewMessage()
	msg.Header.Set(field.NewMsgType(enum.MsgType_ORDER_SINGLE))

	msg.Body.Set(field.NewClOrdID(id.String()))
	msg.Body.Set(field.NewSymbol(instrument))
	msg.Body.Set(field.NewSide(side))
	msg.Body.SetString(tag.OrderQty, floatToStr(amount))
	msg.Body.SetString(tag.Price, floatToStr(price))
	msg.Body.Set(field.NewOrdType(orderType))
	msg.Body.Set(field.NewTimeInForce(timeInForce))
	if execInst != "" {
		msg.Body.SetString(tag.ExecInst, execInst)
	}
	if label != "" {
		msg.Body.SetString(tagDeribitLabel, label)
	}

	resp, err := c.Call(ctx, id.String(), msg)
	if err != nil {
		c.log.Errorw(
			"Fail to create new order",
			"request", msg,
			"error", err,
		)
		return order, err
	}

	order, err = decodeExecutionReport(resp)
	if err != nil {
		c.log.Errorw(
			"Fail to decode ExecutionReport message",
			"request", msg,
			"response", resp,
			"error", err,
		)
		return order, err
	}

	order.TimeInForce = decodeTimeInForce(timeInForce)
	order.OriginalOrderType = decodeOrderType(orderType)
	if strings.Contains(execInst, string(enum.ExecInst_PARTICIPANT_DONT_INITIATE)) {
		order.PostOnly = true
	}
	if strings.Contains(execInst, string(enum.ExecInst_DO_NOT_INCREASE)) {
		order.ReduceOnly = true
	}

	return order, nil
}
