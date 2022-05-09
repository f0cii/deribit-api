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
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

const nonceLen = 64

// Client implements the quickfix.Application interface.
type Client struct {
	log *zap.SugaredLogger

	apiKey    string
	secretKey string

	settings *quickfix.Settings

	targetCompID string
	senderCompID string

	initiator *quickfix.Initiator

	mu          sync.Mutex
	isConnected bool

	currentID int64

	sending sync.Mutex
	pending map[string]*call

	subscriptions    []string
	subscriptionsMap map[string]bool
	emitter          *emission.Emitter
}

// OnCreate implemented as part of Application interface.
func (c *Client) OnCreate(_ quickfix.SessionID) {}

// OnLogon implemented as part of Application interface.
func (c *Client) OnLogon(_ quickfix.SessionID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.isConnected = true

	c.log.Debugw("Logon successfully!")

	if len(c.subscriptions) > 0 {
		err := c.Subscribe(context.Background(), c.subscriptions)
		if err != nil {
			c.log.Warnw("Fail to resubscribe to channels", "error", err)
		}
	}
}

// OnLogout implemented as part of Application interface.
func (c *Client) OnLogout(_ quickfix.SessionID) {
	c.mu.Lock()
	c.isConnected = false
	c.mu.Unlock()

	c.log.Debugw("Logged out!")
	c.sending.Lock()
	for _, call := range c.pending {
		call.done <- ErrClosed
		close(call.done)
	}
	c.sending.Unlock()
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

	var reqIDTag quickfix.Tag
	switch enum.MsgType(msgType) {
	case enum.MsgType_SECURITY_LIST:
		reqIDTag = tag.SecurityReqID
	case enum.MsgType_MARKET_DATA_REQUEST_REJECT:
		reqIDTag = tag.MDReqID
	case enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH:
		reqIDTag = tag.MDReqID
	case enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH:
		reqIDTag = tag.MDReqID
	case enum.MsgType_EXECUTION_REPORT:
		reqIDTag = tag.OrigClOrdID
	case enum.MsgType_ORDER_CANCEL_REJECT:
		reqIDTag = tag.ClOrdID
	case enum.MsgType_ORDER_MASS_CANCEL_REPORT:
		reqIDTag = tag.OrderID
	case enum.MsgType_POSITION_REPORT:
		reqIDTag = tag.PosReqID
	case enum.MsgType_USER_RESPONSE:
		reqIDTag = tag.UserRequestID
	case enum.MsgType_SECURITY_STATUS:
		reqIDTag = tag.SecurityStatusReqID
	}

	id, err := msg.Body.GetString(reqIDTag)
	if err != nil {
		c.log.Errorw("Fail to get request ID", "error", err)
		return err
	}

	c.mu.Lock()
	call := c.pending[id]
	delete(c.pending, id)
	c.mu.Unlock()

	if call != nil {
		call.response = msg
	}

	if call != nil {
		call.done <- nil
		close(call.done)
	}

	return nil
}

// New returns a new client for Deribit FIX API.
func New(
	ctx context.Context,
	apiKey string,
	secretKey string,
	settings *quickfix.Settings,
) (*Client, error) {
	l := zap.S()

	// Get TargetCompID and SenderCompID from settings.
	globalSettings := settings.GlobalSettings()
	targetCompID, err := globalSettings.Setting("TargetCompID")
	if err != nil {
		l.Errorw("Fail to read TargetCompID from settings", "error", err)
		return nil, err
	}

	senderCompID, err := globalSettings.Setting("SenderCompID")
	if err != nil {
		l.Errorw("Fail to read SenderCompID from settings", "error", err)
		return nil, err
	}

	// Create a new Client object.
	client := &Client{
		log:              l,
		apiKey:           apiKey,
		secretKey:        secretKey,
		settings:         settings,
		targetCompID:     targetCompID,
		senderCompID:     senderCompID,
		mu:               sync.Mutex{},
		isConnected:      false,
		sending:          sync.Mutex{},
		pending:          make(map[string]*call),
		subscriptionsMap: make(map[string]bool),
		emitter:          emission.NewEmitter(),
	}

	// Init session and logon to deribit FIX API server.
	logFactory := quickfix.NewNullLogFactory()

	client.initiator, err = quickfix.NewInitiator(
		client,
		quickfix.NewMemoryStoreFactory(),
		settings,
		logFactory,
	)
	if err != nil {
		client.log.Errorw("Fail to create new initiator", "error", err)
		return nil, err
	}

	if err = client.initiator.Start(); err != nil {
		client.log.Errorw("Fail to initialize initiator", "error", err)
		return nil, err
	}

	// Wait for the session to be authorized by the server.
	for !client.IsConnected() {
		time.Sleep(10 * time.Millisecond)
	}

	return client, nil
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

func (c *Client) getCurrentID() int64 {
	c.currentID++
	return c.currentID
}

func (c *Client) handleSubscriptions(msgType string, msg *quickfix.Message) {
	switch enum.MsgType(msgType) {
	case enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH, enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH:
		symbol, err := getSymbol(msg)
		if err != nil {
			c.log.Warnw("Fail to get symbol", "error", err)
			return
		}

		entries, err := getMDEntries(msg)
		if err != nil {
			c.log.Warnw(
				"Fail to get NoMDEntries",
				"msg", msg,
				"error", err,
			)
			return
		}

		sendTime, err := getSendingTime(msg)
		if err != nil {
			c.log.Warnw("Fail to get SendingTime", "error", err)
			return
		}

		orderBookEvent := models.OrderBookRawNotification{
			Timestamp:      sendTime.UnixMilli(),
			InstrumentName: symbol,
		}
		for i := 0; i < entries.Len(); i++ {
			entry := entries.Get(i)
			entryType, err := getMDEntryType(entry)
			if err != nil {
				c.log.Debugw("No value for MDEntryType", "error", err)
				continue
			}

			if entryType != enum.MDEntryType_BID &&
				entryType != enum.MDEntryType_OFFER {
				continue
			}

			price, err := getMDEntryPx(entry)
			if err != nil {
				c.log.Debugw("No value for MDEntryPx", "error", err)
				continue
			}

			action, err := getMDUpdateAction(entry)
			if err != nil {
				action = "new"
			}

			var amount decimal.Decimal
			if action != "delete" {
				amount, err = getMDEntrySize(entry)
				if err != nil {
					c.log.Debugw("No value for MDEntrySize", "error", err)
					continue
				}
			}

			item := models.OrderBookNotificationItem{
				Action: action,
				Price:  price,
				Amount: amount,
			}
			if entryType == enum.MDEntryType_BID {
				orderBookEvent.Bids = append(orderBookEvent.Bids, item)
			} else {
				orderBookEvent.Asks = append(orderBookEvent.Asks, item)
			}
		}

		if len(orderBookEvent.Bids) > 0 || len(orderBookEvent.Asks) > 0 {
			var isSnapshot bool
			if msgType == string(enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH) {
				isSnapshot = true
			}

			c.Emit(newOrderBookNotificationChannel(symbol), orderBookEvent, isSnapshot)
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
	msg.Body.Set(field.NewMDReqID(id))

	var cc *call
	if wait {
		cc = &call{request: msg, done: make(chan error, 1)}
		c.pending[id] = cc
	}
	c.mu.Unlock()

	if err := quickfix.Send(msg); err != nil {
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
	marketDepth int,
	mdUpdateType enum.MDUpdateType,
	mdEntryTypes []enum.MDEntryType,
	instruments []string,
) (*quickfix.Message, error) {
	id := strconv.FormatInt(c.getCurrentID(), 10)

	msg := quickfix.NewMessage()
	msg.Header.Set(field.NewMsgType(enum.MsgType_MARKET_DATA_REQUEST))

	msg.Body.Set(field.NewSubscriptionRequestType(subscriptionRequestType))
	msg.Body.Set(field.NewMarketDepth(marketDepth))
	if subscriptionRequestType == enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES {
		msg.Body.Set(field.NewMDUpdateType(mdUpdateType))
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

	return c.Call(ctx, id, msg)
}

func (c *Client) SubscribeOrderBooks(ctx context.Context, instruments []string) error {
	if len(instruments) == 0 {
		c.log.Debugw("No instruments to subscribe")
		return nil
	}

	msg, err := c.MarketDataRequest(
		ctx,
		enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
		0,
		enum.MDUpdateType_INCREMENTAL_REFRESH,
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

// Subscribe listens for notifications.
// Currently, only support for orderbook notifications.
func (c *Client) Subscribe(ctx context.Context, channels []string) error {
	c.mu.Lock()
	instruments := make([]string, 0, len(channels))
	for _, channel := range channels {
		if !c.subscriptionsMap[channel] {
			parts := strings.Split(channel, ".")
			if len(parts) == 2 && parts[0] == "book" {
				c.subscriptionsMap[channel] = true
				c.subscriptions = append(c.subscriptions, channel)
				instruments = append(instruments, parts[1])
			}
		}
	}
	c.mu.Unlock()

	return c.SubscribeOrderBooks(ctx, instruments)
}

func (c *Client) CreateOrder(
	ctx context.Context,
	instrument string,
	side enum.Side,
	amount decimal.Decimal,
	price decimal.Decimal,
	orderType enum.OrdType,
	timeInForce enum.TimeInForce,
	execInst string,
	label string,
) (order models.Order, err error) {
	id := strconv.FormatInt(c.getCurrentID(), 10)

	msg := quickfix.NewMessage()
	msg.Header.Set(field.NewMsgType(enum.MsgType_ORDER_SINGLE))

	msg.Body.Set(field.NewClOrdID(id))
	msg.Body.Set(field.NewSymbol(instrument))
	msg.Body.Set(field.NewSide(side))
	msg.Body.Set(field.NewOrderQty(amount, 0))
	msg.Body.Set(field.NewPrice(price, 0))
	msg.Body.Set(field.NewOrdType(orderType))
	msg.Body.Set(field.NewTimeInForce(timeInForce))
	if execInst != "" {
		msg.Body.SetString(tag.ExecInst, execInst)
	}
	if label != "" {
		msg.Body.SetString(tagDeribitLabel, label)
	}

	resp, err := c.Call(ctx, id, msg)
	if err != nil {
		c.log.Errorw("Fail to create new order", "error", err)
		return
	}

	order, err = decodeExecutionReport(resp)
	if err != nil {
		c.log.Errorw("Fail to decode ExecutionReport message", "error", err)
		return
	}

	order.TimeInForce = decodeTimeInForce(timeInForce)
	order.OriginalOrderType = decodeOrderType(orderType)
	if strings.Contains(execInst, string(enum.ExecInst_PARTICIPANT_DONT_INITIATE)) {
		order.PostOnly = true
	}
	if strings.Contains(execInst, string(enum.ExecInst_DO_NOT_INCREASE)) {
		order.ReduceOnly = true
	}

	return
}
