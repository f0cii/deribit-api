package deribit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/deribit-api/models"
	"github.com/chuckpreslar/emission"
	ws "github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	sws "github.com/sourcegraph/jsonrpc2/websocket"
	"go.uber.org/zap"
)

const (
	RealBaseURL = "wss://www.deribit.com/ws/api/v2/"
	TestBaseURL = "wss://test.deribit.com/ws/api/v2/"
)

var (
	ErrAuthenticationIsRequired = errors.New("authentication is required")
	ErrNotConnected             = errors.New("not connected")
)

// Event is wrapper of received event
type Event struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type Configuration struct {
	Addr          string `json:"addr"`
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
}

type Client struct {
	l *zap.SugaredLogger

	addr          string
	apiKey        string
	secretKey     string
	autoReconnect bool
	debugMode     bool

	conn        *ws.Conn
	rpcConn     *jsonrpc2.Conn
	mu          sync.RWMutex
	once        sync.Once
	heartCancel chan struct{}
	isConnected bool

	subscriptions    []string
	subscriptionsMap map[string]struct{}

	emitter *emission.Emitter
}

func New(l *zap.SugaredLogger, cfg *Configuration) *Client {
	return &Client{
		l:                l,
		addr:             cfg.Addr,
		apiKey:           cfg.ApiKey,
		secretKey:        cfg.SecretKey,
		autoReconnect:    cfg.AutoReconnect,
		debugMode:        cfg.DebugMode,
		mu:               sync.RWMutex{},
		once:             sync.Once{},
		subscriptionsMap: make(map[string]struct{}),
		emitter:          emission.NewEmitter(),
	}
}

// setIsConnected sets state for isConnected
func (c *Client) setIsConnected(state bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isConnected = state
}

// IsConnected returns the WebSocket connection state
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.isConnected
}

func (c *Client) Subscribe(channels []string) error {
	return c.subscribe(channels, true)
}

func (c *Client) subscribe(channels []string, isNewSubscription bool) error {
	l := c.l.With("func", "subscribe")
	var publicChannels []string
	var privateChannels []string

	for _, v := range channels {
		if _, ok := c.subscriptionsMap[v]; ok {
			continue
		}
		if strings.HasPrefix(v, "user.") {
			privateChannels = append(privateChannels, v)
		} else {
			publicChannels = append(publicChannels, v)
		}
	}

	if len(publicChannels) > 0 {
		pubSubResp, err := c.PublicSubscribe(context.Background(), &models.SubscribeParams{
			Channels: publicChannels,
		})
		if err != nil {
			l.Errorw("error subscribe public", "err", err)
			return err
		}
		if isNewSubscription {
			c.subscriptions = append(c.subscriptions, pubSubResp...)
		}
		for _, v := range pubSubResp {
			c.subscriptionsMap[v] = struct{}{}
		}
	}

	if len(privateChannels) > 0 {
		privateSubResp, err := c.PrivateSubscribe(context.Background(), &models.SubscribeParams{
			Channels: privateChannels,
		})
		if err != nil {
			l.Errorw("error subscribe private", "err", err)
			return err
		}
		if isNewSubscription {
			c.subscriptions = append(c.subscriptions, privateSubResp...)
		}
		for _, v := range privateSubResp {
			c.subscriptionsMap[v] = struct{}{}
		}
	}
	return nil
}

// Start connect ws
func (c *Client) Start() error {
	c.setIsConnected(false)
	c.subscriptionsMap = make(map[string]struct{})
	c.conn = nil
	c.rpcConn = nil
	c.heartCancel = make(chan struct{})

	var (
		err  error
		conn *ws.Conn
	)
	for i := 0; i < 3; i++ {
		conn, _, err = ws.DefaultDialer.Dial(c.addr, nil)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		c.conn = conn
		break
	}
	if err != nil {
		return err
	}

	c.rpcConn = jsonrpc2.NewConn(context.Background(), sws.NewObjectStream(c.conn), c)

	c.setIsConnected(true)

	// auth
	if c.apiKey != "" && c.secretKey != "" {
		if _, err := c.Auth(context.Background()); err != nil {
			return fmt.Errorf("failed to auth, err = %s", err)
		}
	}

	// subscribe
	if err := c.subscribe(c.subscriptions, false); err != nil {
		return fmt.Errorf("failed to subscribe, err=%s", err)
	}

	if _, err := c.SetHeartbeat(context.Background(), &models.SetHeartbeatParams{Interval: 30}); err != nil {
		return fmt.Errorf("failed to set heartbeat, err=%s", err)
	}

	go c.heartbeat()

	c.once.Do(func() {
		if c.autoReconnect {
			go c.reconnect()
		}
	})

	return nil
}

// Call issues JSONRPC v2 calls
func (c *Client) Call(ctx context.Context, method string, params interface{}, result interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if !c.IsConnected() {
		return ErrNotConnected
	}
	if params == nil {
		params = json.RawMessage("{}")
	}

	return c.rpcConn.Call(ctx, method, params, result)
}

// Handle implements jsonrpc2.Handler
func (c *Client) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	if req.Method == "subscription" {
		if req.Params != nil && len(*req.Params) > 0 {
			var event Event
			if err := json.Unmarshal(*req.Params, &event); err != nil {
				return
			}
			c.subscriptionsProcess(&event)
		}
	}
}

// ResetConnection force reconnect
func (c *Client) ResetConnection() {
	_ = c.conn.Close()
}

func (c *Client) heartbeat() {
	l := c.l.With("func", "heartbeat")
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			if _, err := c.Test(context.Background()); err != nil {
				l.Errorw("error test server", "err", err)
				_ = c.conn.Close() // close server
			}
		case <-c.heartCancel:
			return
		}
	}
}

func (c *Client) reconnect() {
	l := c.l.With("func", "reconnect")
	for {
		notify := c.rpcConn.DisconnectNotify()
		<-notify
		c.setIsConnected(false)

		l.Infow("disconnect, reconnect...")

		close(c.heartCancel)

		time.Sleep(1 * time.Second)

		for {
			if err := c.Start(); err != nil {
				if c.rpcConn != nil {
					_ = c.rpcConn.Close()
				}
				l.Errorw("reconnect: start error", "err", err)
				time.Sleep(5 * time.Second)
			} else {
				l.Infow("reconnect successfully")
				break
			}
		}
	}
}
