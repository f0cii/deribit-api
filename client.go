package deribit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/KyberNetwork/deribit-api/models"
	"github.com/chuckpreslar/emission"
	ws "github.com/gorilla/websocket"
	"github.com/sourcegraph/jsonrpc2"
	sws "github.com/sourcegraph/jsonrpc2/websocket"
)

const (
	RealBaseURL = "wss://www.deribit.com/ws/api/v2/"
	TestBaseURL = "wss://test.deribit.com/ws/api/v2/"
)

const (
	MaxTryTimes = 10
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
	addr          string
	apiKey        string
	secretKey     string
	autoReconnect bool
	debugMode     bool

	conn        *ws.Conn
	rpcConn     *jsonrpc2.Conn
	mu          sync.RWMutex
	heartCancel chan struct{}
	isConnected bool

	subscriptions    []string
	subscriptionsMap map[string]struct{}

	emitter *emission.Emitter
}

func New(cfg *Configuration) (*Client, error) {
	client := &Client{
		addr:             cfg.Addr,
		apiKey:           cfg.ApiKey,
		secretKey:        cfg.SecretKey,
		autoReconnect:    cfg.AutoReconnect,
		debugMode:        cfg.DebugMode,
		subscriptionsMap: make(map[string]struct{}),
		emitter:          emission.NewEmitter(),
	}
	err := client.start()
	if err != nil {
		return nil, err
	}
	return client, nil
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

func (c *Client) Subscribe(channels []string) {
	c.subscriptions = append(c.subscriptions, channels...)
	c.subscribe(channels)
}

func (c *Client) subscribe(channels []string) {
	var publicChannels []string
	var privateChannels []string

	for _, v := range c.subscriptions {
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
		if _, err := c.PublicSubscribe(context.Background(), &models.SubscribeParams{
			Channels: publicChannels,
		}); err != nil {
			log.Printf("error subscribe public err = %s", err)
		}
	}
	if len(privateChannels) > 0 {
		if _, err := c.PrivateSubscribe(context.Background(), &models.SubscribeParams{
			Channels: privateChannels,
		}); err != nil {
			log.Printf("error subscribe private err = %s", err)
		}
	}

	allChannels := append(publicChannels, privateChannels...)
	for _, v := range allChannels {
		c.subscriptionsMap[v] = struct{}{}
	}
}

func (c *Client) start() error {
	c.setIsConnected(false)
	c.subscriptionsMap = make(map[string]struct{})
	c.conn = nil
	c.rpcConn = nil
	c.heartCancel = make(chan struct{})

	var (
		err  error
		conn *ws.Conn
	)
	for i := 0; i < MaxTryTimes; i++ {
		conn, _, err = ws.DefaultDialer.Dial(c.addr, nil)
		if err != nil {
			time.Sleep(time.Duration(i+1) * time.Second)
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
			log.Printf("failed to auth, err = %s", err)
			return err
		}
	}

	// subscribe
	c.subscribe(c.subscriptions)

	if _, err := c.SetHeartbeat(context.Background(), &models.SetHeartbeatParams{Interval: 30}); err != nil {
		return err
	}

	if c.autoReconnect {
		go c.reconnect()
	}

	go c.heartbeat()

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

func (c *Client) heartbeat() {
	t := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-t.C:
			if _, err := c.Test(context.Background()); err != nil {
				log.Printf("error test server, err = %s", err)
				_ = c.conn.Close() // close server
			}
		case <-c.heartCancel:
			return
		}
	}
}

func (c *Client) reconnect() {
	notify := c.rpcConn.DisconnectNotify()
	<-notify
	c.setIsConnected(false)

	log.Println("disconnect, reconnect...")

	close(c.heartCancel)

	time.Sleep(1 * time.Second)

	_ = c.start()
}
