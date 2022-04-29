package fix

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"sync"
	"time"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

const nonceLen = 64

// Client implements the quickfix.Application interface.
type Client struct {
	l *zap.SugeredLogger

	addr      string
	apiKey    string
	secretKey string

	settings *quickfix.Settings

	mu          sync.Mutex
	isConnected bool
}

// OnCreate implemented as part of Application interface.
func (c *Client) OnCreate(_ quickfix.SessionID) {}

// OnLogon implemented as part of Application interface.
func (c *Client) OnLogon(_ quickfix.SessionID) {
	mu.Lock()
	defer mu.Unlock()
	isConnected = true
}

// OnLogout implemented as part of Application interface.
func (c *Client) OnLogout(_ quickfix.SessionID) {
	mu.Lock()
	defer mu.Unlock()
	isConnected = false
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
		log.Errorw(
			"Fail to generate random bytes",
			"nonce_len", c.nonceLen,
			"error", err,
		)
		return
	}

	rawData := strconv.FormatInt(timestamp, 10) + "." + base64.StdEncoding.EncodeToString(nonce)
	hash := sha256.Sum256([]byte(rawData + secretKey))
	password := base64.StdEncoding.EncodeToString(hash[:])

	msg.Body.SetField(tag.RawData, quickfix.FIXString(rawData))
	msg.Body.SetField(tag.Username, quickfix.FIXString(apiKey))
	msg.Body.SetField(tag.Password, quickfix.FIXString(password))
}

// ToApp implemented as a part of Application interface.
func (c *Client) ToApp(msg *quickfix.Message, _ quickfix.SessionID) error {
	c.log.Debugw("Sending message to server", "msg", msg)
}

// FromApp implemented as a part of Application interface.
func (c *Client) FromApp(msg *quickfix.Message, _ quickfix.SessionID) quickfix.MessageRejectError {
	// Process message according to message type.
}

// New returns a new client for Deribit FIX API.
func New(addr, apiKey, secretKey string, settings *quickfix.Settings) (*Client, error) {
	// Create a new Client object.
	client := &Client{
		l:           zap.S(),
		addr:        addr,
		apiKey:      apiKey,
		secretKey:   secretKey,
		settings:    settings,
		mu:          sync.Mutex{},
		isConnected: false,
	}

	// Init session and logon to deribit FIX API server.

	// Waiting for the response.

	return client, nil
}
