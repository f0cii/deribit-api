package deribit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chuckpreslar/emission"
	"github.com/sourcegraph/jsonrpc2"
	"github.com/sumorf/deribit-api/models"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"strings"
	"sync"
	"time"
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
)

// Event is wrapper of received event
type Event struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

type Configuration struct {
	Ctx           context.Context
	Addr          string `json:"addr"`
	ApiKey        string `json:"api_key"`
	SecretKey     string `json:"secret_key"`
	AutoReconnect bool   `json:"auto_reconnect"`
	DebugMode     bool   `json:"debug_mode"`
}

type Client struct {
	ctx           context.Context
	addr          string
	apiKey        string
	secretKey     string
	autoReconnect bool

	conn        *websocket.Conn
	rpcConn     *jsonrpc2.Conn
	mu          sync.RWMutex
	heartCancel chan struct{}
	isConnected bool

	auth struct {
		token   string
		refresh string
	}

	subscriptions    []string
	subscriptionsMap map[string]struct{}

	emitter *emission.Emitter
}

func New(cfg *Configuration) *Client {
	ctx := cfg.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	client := &Client{
		ctx:              ctx,
		addr:             cfg.Addr,
		apiKey:           cfg.ApiKey,
		secretKey:        cfg.SecretKey,
		autoReconnect:    cfg.AutoReconnect,
		subscriptionsMap: make(map[string]struct{}),
		emitter:          emission.NewEmitter(),
	}
	err := client.start()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// setIsConnected sets state for isConnected
func (b *Client) setIsConnected(state bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.isConnected = state
}

// IsConnected returns the WebSocket connection state
func (b *Client) IsConnected() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.isConnected
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
		c.PublicSubscribe(&models.SubscribeParams{
			Channels: publicChannels,
		})
	}
	if len(privateChannels) > 0 {
		c.PrivateSubscribe(&models.SubscribeParams{
			Channels: privateChannels,
		})
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

	for i := 0; i < MaxTryTimes; i++ {
		conn, _, err := c.connect()
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		c.conn = conn
		break
	}
	if c.conn == nil {
		return errors.New("connect fail")
	}

	c.rpcConn = jsonrpc2.NewConn(context.Background(), NewObjectStream(c.conn), c)

	// auth
	if c.apiKey != "" && c.secretKey != "" {
		if err := c.Auth(c.apiKey, c.secretKey); err != nil {
			log.Printf("auth error: %v", err)
		}
	}

	// subscribe
	c.subscribe(c.subscriptions)

	c.SetHeartbeat(&models.SetHeartbeatParams{Interval: 30})

	if c.autoReconnect {
		go c.reconnect()
	}

	c.setIsConnected(true)

	go c.heartbeat()

	return nil
}

// Call issues JSONRPC v2 calls
func (c *Client) Call(method string, params interface{}, result interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("%v", r))
		}
	}()

	if !c.IsConnected() {
		return errors.New("not connected")
	}
	if params == nil {
		params = emptyParams
	}

	if token, ok := params.(privateParams); ok {
		if c.auth.token == "" {
			return ErrAuthenticationIsRequired
		}
		token.setToken(c.auth.token)
	}

	return c.rpcConn.Call(c.ctx, method, params, result)
}

// Handle implements jsonrpc2.Handler
func (c *Client) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	//log.Printf("Handle %v", req.Method)
	if req.Method == "subscription" {
		// update events
		if req.Params != nil && len(*req.Params) > 0 {
			var event Event
			if err := json.Unmarshal(*req.Params, &event); err != nil {
				//c.setError(err)
				return
			}
			c.subscriptionsProcess(&event)
		}
	}
}

func (c *Client) heartbeat() {
	t := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-t.C:
			c.Test()
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

	c.start()
}

func (c *Client) connect() (*websocket.Conn, *http.Response, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	conn, resp, err := websocket.Dial(ctx, c.addr, websocket.DialOptions{})
	return conn, resp, err
}
