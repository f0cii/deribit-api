package fix

import (
	"bytes"
	"errors"
	"sync"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

type MockInitiator struct {
	mu              *sync.Mutex
	app             quickfix.Application
	settings        *quickfix.Settings
	sessionSettings map[quickfix.SessionID]*quickfix.SessionSettings
	storeFactory    quickfix.MessageStoreFactory
	logFactory      quickfix.LogFactory
	results         []interface{}
	isActive        bool
}

func createMockInitiator(
	app quickfix.Application,
	storeFactory quickfix.MessageStoreFactory,
	appSettings *quickfix.Settings,
	logFactory quickfix.LogFactory,
) (Initiator, error) {
	i := &MockInitiator{
		mu:              &sync.Mutex{},
		app:             app,
		storeFactory:    storeFactory,
		settings:        appSettings,
		sessionSettings: appSettings.SessionSettings(),
		logFactory:      logFactory,
		results: []interface{}{
			&models.AuthResponse{},
			"success",
		},
	}
	for sessionID := range i.sessionSettings {
		app.OnCreate(sessionID)
	}

	return i, nil
}

func (i *MockInitiator) Start() error {
	for sessionID, s := range i.sessionSettings {
		// send Logon message
		i.app.ToAdmin(quickfix.NewMessage(), sessionID)

		if !s.HasSetting("SocketConnectHost") {
			return errors.New("Conditionally Required Setting: SocketConnectHost")
		}

		if !s.HasSetting("SocketConnectPort") {
			return errors.New("Conditionally Required Setting: SocketConnectPort")
		}
		i.app.OnLogon(sessionID)
	}

	i.setIsActive(true)
	return nil
}

func (i *MockInitiator) Stop() {
	if i.getIsActive() {
		for sessionID := range i.sessionSettings {
			i.app.OnLogout(sessionID)
		}
		i.setIsActive(false)
		return
	}
}

func (i *MockInitiator) getIsActive() bool {
	i.mu.Lock()
	defer i.mu.Unlock()

	return i.isActive
}

func (i *MockInitiator) setIsActive(val bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	i.isActive = val
}

// Send sends message to counterparty (Deribit Server)
func (i *MockInitiator) send(msg *quickfix.Message) error {
	msgType, err := msg.Header.GetBytes(tag.MsgType)
	if err != nil {
		return err
	}

	var err2 error
	var reqIDTag quickfix.Tag
	fixMsgType := enum.MsgType(msgType)
	if fixMsgType == enum.MsgType_ORDER_SINGLE {
		reqIDTag = tag.ClOrdID
	} else {
		reqIDTag, err2 = getReqIDTagFromMsgType(enum.MsgType(msgType))
		if err2 != nil {
			return err2
		}
	}

	mutex.Lock()
	requestID, err = msg.Body.GetString(reqIDTag)
	mutex.Unlock()
	if err != nil {
		return err
	}

	if isAdminMessageType(msgType) {
		for sessionID := range i.sessionSettings {
			i.app.ToAdmin(msg, sessionID)
		}
	} else {
		for sessionID := range i.sessionSettings {
			err := i.app.ToApp(msg, sessionID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Receive receives message from counterparty (Deribit Server)
func (i *MockInitiator) receive(msg *quickfix.Message) error {
	msgType, err := msg.Header.GetBytes(tag.MsgType)
	if err != nil {
		return err
	}

	if isAdminMessageType(msgType) {
		for sessionID := range i.sessionSettings {
			err := i.app.FromAdmin(msg, sessionID)
			if err != nil {
				return err
			}
		}
	} else {
		for sessionID := range i.sessionSettings {
			err := i.app.FromApp(msg, sessionID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func mockDeribitResponse(msg *quickfix.Message) (err error) {
	initiator := mockInitiator.(*MockInitiator)
	return initiator.receive(msg)
}

func mockSender(m quickfix.Messagable) (err error) {
	initiator := mockInitiator.(*MockInitiator)
	return initiator.send(m.ToMessage())
}

// nolint:gochecknoglobals
var (
	msgTypeHeartbeat     = []byte("0")
	msgTypeLogon         = []byte("A")
	msgTypeTestRequest   = []byte("1")
	msgTypeResendRequest = []byte("2")
	msgTypeReject        = []byte("3")
	msgTypeSequenceReset = []byte("4")
	msgTypeLogout        = []byte("5")
)

// isAdminMessageType returns true if the message type is a session level message.
func isAdminMessageType(m []byte) bool {
	switch {
	case bytes.Equal(msgTypeHeartbeat, m),
		bytes.Equal(msgTypeLogon, m),
		bytes.Equal(msgTypeTestRequest, m),
		bytes.Equal(msgTypeResendRequest, m),
		bytes.Equal(msgTypeReject, m),
		bytes.Equal(msgTypeSequenceReset, m),
		bytes.Equal(msgTypeLogout, m):
		return true
	}

	return false
}
