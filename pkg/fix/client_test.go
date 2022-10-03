package fix

import (
	"bytes"
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
	"github.com/stretchr/testify/suite"
)

const (
	apiKey       = "api_key"
	secretKey    = "secret_key"
	responseTime = 100 * time.Microsecond
)

// nolint:gochecknoglobals
var (
	mockInitiator Initiator
	requestID     string
	mutex         = sync.Mutex{}
)

type FixTestSuite struct {
	suite.Suite
	c *Client
}

func TestFixTestSuite(t *testing.T) {
	suite.Run(t, new(FixTestSuite))
}

// nolint:lll
func (ts *FixTestSuite) SetupSuite() {
	require := ts.Require()

	initiateFixClientTests := []struct {
		config        string
		requiredError bool
	}{
		{
			"[DEFAULT]\nSocketConnectHost=test.deribit.com\nSocketConnectPort=9881\nHeartBtInt=30\nSenderCompID=FIX_TEST\nTargetCompID=DERIBITSERVER\nResetOnLogon=Y\n\n[SESSION]\nBeginString=FIX.4.4\n",
			false,
		},
		{
			"[DEFAULT]\nSocketConnectPort=9881\nHeartBtInt=30\nSenderCompID=FIX_TEST\nTargetCompID=DERIBITSERVER\nResetOnLogon=Y\n\n[SESSION]\nBeginString=FIX.4.4\n",
			true, //  "Conditionally Required Setting: SocketConnectHost"
		},
		{
			"[DEFAULT]\nSocketConnectHost=test.deribit.com\nSocketConnectPort=9881\nHeartBtInt=30\nSenderCompID=FIX_TEST\nResetOnLogon=Y\n\n[SESSION]\nBeginString=FIX.4.4\n",
			true, //  "Conditionally Required Setting: TargetCompID"
		},
		{
			"[DEFAULT]\nSocketConnectHost=test.deribit.com\nSocketConnectPort=9881\nHeartBtInt=30\nTargetCompID=DERIBITSERVER\nResetOnLogon=Y\n\n[SESSION]\nBeginString=FIX.4.4\n",
			true, //  "Conditionally Required Setting: SenderCompID"
		},
	}

	for _, test := range initiateFixClientTests {
		appSettings, err := quickfix.ParseSettings(bytes.NewBufferString(test.config))
		require.NoError(err)
		cfg := Config{
			APIKey:    apiKey,
			SecretKey: secretKey,
			Settings:  appSettings,
			Dialer:    createMockInitiator,
			Sender:    mockSender,
		}
		c, err := New(context.Background(), cfg)
		if test.requiredError {
			require.Error(err)
		} else {
			require.NoError(err)
			mockInitiator = c.initiator
			ts.c = c
		}
	}

	// Default Initiator and MockSender
	settingStr := "[DEFAULT]\nSenderCompID=FIX_TEST\nTargetCompID=XXX\nResetOnLogon=Y\n\n[SESSION]\nBeginString=FIX.4.4\n"
	appSettings, err := quickfix.ParseSettings(bytes.NewBufferString(settingStr))
	require.NoError(err)
	cfg := Config{
		APIKey:    apiKey,
		SecretKey: secretKey,
		Settings:  appSettings,
	}

	client, err := New(context.Background(), cfg)
	require.Error(err)
	require.Nil(client)
}

// nolint:lll,funlen
func (ts *FixTestSuite) TestHandleSubscriptions() {
	require := ts.Require()

	type orderbookEvent struct {
		event *models.OrderBookRawNotification
		reset bool
	}

	tests := []struct {
		msgType      string
		fixMsg       string
		channel      string
		expectOutput interface{}
	}{
		{
			"W", // enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH
			"8=FIX.4.4\u00019=293\u000135=W\u000149=DERIBITSERVER\u000156=OPTION_TRADING_BTC_TESTNET\u000134=2\u000152=20220815-10:39:22.035\u000155=BTC-26AUG22-32000-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=24185.9900\u0001100087=0.0000\u0001100090=0.3238\u0001746=0.0000\u0001201=0\u0001262=8cd489c3-1045-4e53-a9e5-7926ec3579c0\u0001268=1\u0001269=1\u0001270=0.8735\u0001271=6.0000\u0001272=20220815-10:39:21.568\u000110=116\u0001",
			"book.BTC-26AUG22-32000-P",
			orderbookEvent{
				&models.OrderBookRawNotification{
					Timestamp:      1660559961568,
					InstrumentName: "BTC-26AUG22-32000-P",
					PrevChangeID:   0,
					ChangeID:       0,
					Asks: []models.OrderBookNotificationItem{
						{
							Action: "new",
							Price:  0.8735,
							Amount: 6,
						},
					},
				},
				true,
			},
		},
		{
			"W", // enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH
			"8=FIX.4.4\u00019=353\u000135=W\u000149=DERIBITSERVER\u000156=OPTION_TRADING_BTC_TESTNET\u000134=2\u000152=20220804-08:54:42.073\u000155=BTC-28OCT22-32000-P\u0001231=1.0000\u0001311=SYN.BTC-28OCT22\u0001810=22943.2054\u0001100087=0.0000\u0001100090=0.4305\u0001746=1.0000\u0001201=0\u0001262=7c268500-604f-45df-a4eb-7954d74e89ab\u0001268=2\u0001269=0\u0001270=0.4005\u0001271=12.0000\u0001272=20220804-08:54:41.698\u0001269=1\u0001270=0.4545\u0001271=12.0000\u0001272=20220804-08:54:41.698\u000110=132\u0001",
			"book.BTC-28OCT22-32000-P",
			orderbookEvent{
				&models.OrderBookRawNotification{
					Timestamp:      1659603281698,
					InstrumentName: "BTC-28OCT22-32000-P",
					PrevChangeID:   0,
					ChangeID:       0,
					Bids: []models.OrderBookNotificationItem{
						{
							Action: "new",
							Price:  0.4005,
							Amount: 12,
						},
					},
					Asks: []models.OrderBookNotificationItem{
						{
							Action: "new",
							Price:  0.4545,
							Amount: 12,
						},
					},
				},
				true,
			},
		},
		{
			"7", // enum.MsgType_ADVERTISEMENT,
			"8=FIX.4.4\u00019=5\u000135=7\u000110=170\u0001",
			"",
			nil,
		},
		// some wrong decoder
		{
			"W",
			"8=FIX.4.4\u00019=5\u000135=7\u000110=170\u0001",
			"",
			nil, // failed to getSymbol
		},
		{
			"W",
			"8=FIX.4.4\u00019=17\u000135=7\u000155=BTC_USDT\u000110=253\u0001",
			"",
			nil, // failed to getSymbol
		},
	}

	eventCh := make(chan interface{}, 100)
	listener := func(event *models.OrderBookRawNotification, reset bool) {
		eventCh <- orderbookEvent{event, reset}
	}

	for _, test := range tests {
		bufferData := bytes.NewBufferString(test.fixMsg)

		msg := quickfix.NewMessage()
		err := quickfix.ParseMessage(msg, bufferData)
		require.NoError(err)

		ts.c.On(test.channel, listener)
		ts.c.handleSubscriptions(test.msgType, msg)
		if test.expectOutput != nil && ts.Len(eventCh, 1) {
			event := <-eventCh
			ts.Assert().Equal(event, test.expectOutput)
		}
		ts.c.Off(test.channel, listener)
	}
}

func (ts *FixTestSuite) TestSend() {
	assert := ts.Assert()
	wait := true

	wrongMsg := quickfix.NewMessage()
	correctMsg := quickfix.NewMessage()
	correctMsg.Body.Set(field.NewOrigClOrdID("test_send_func_0"))
	correctMsg.Header.Set(field.NewMsgType(enum.MsgType_EXECUTION_REPORT))

	tests := []struct {
		msg           *quickfix.Message
		expectOutput  Waiter
		requiredError bool
	}{
		{
			correctMsg,
			Waiter{
				call: &call{
					request: correctMsg,
					done:    make(chan error, 1),
				},
			},
			false,
		},
		{
			wrongMsg,
			Waiter{},
			true, // Conditionally Required Field Missing (35)
		},
	}

	for idx, test := range tests {
		id := "test_send_func_0" + strconv.Itoa(idx)
		waiter, err := ts.c.send(context.Background(), id, test.msg, wait)
		if test.requiredError {
			assert.Error(err)
			assert.Nil(ts.c.pending[id])
		} else {
			assert.NoError(err)
			assert.Equal(test.expectOutput.call.request, waiter.call.request)
			assert.Len(waiter.call.done, 0)
			delete(ts.c.pending, id)
		}
	}
}

func (ts *FixTestSuite) TestCall() {
	assert := ts.Assert()
	require := ts.Require()

	tests := []struct {
		requestMsg    string
		responseMsg   string
		requiredError bool
	}{
		{
			"",
			"",
			true, // send err: Conditionally Required Field Missing (35)
		},
		{
			"8=FIX.4.4\u00019=25\u000135=8\u000141=test_call_func_1\u000110=130\u0001",
			"8=FIX.4.4\u00019=43\u000135=8\u000114=123.4560000000\u000141=test_call_func_1\u000110=204\u0001",
			false,
		},
	}

	for idx, test := range tests {
		id := "test_call_func_" + strconv.Itoa(idx)
		reqMsg := getMsgFromString(test.requestMsg)
		respMsg := getMsgFromString(test.responseMsg)

		if !test.requiredError {
			go func() {
				time.Sleep(responseTime)
				err := mockDeribitResponse(respMsg)
				require.NoError(err)
			}()
		}

		msg, err := ts.c.Call(context.Background(), id, reqMsg)
		assert.Nil(ts.c.pending[id])
		if test.requiredError {
			assert.Error(err)
		} else {
			assert.NoError(err)
			assert.Equal(respMsg.String(), msg.String())
		}
	}
}

// example: request for subscribing orderbook
// nolint:lll
func (ts *FixTestSuite) TestMarketDataRequest() {
	require := ts.Require()
	marketDepth := 0
	mdUpdateType := enum.MDUpdateType_INCREMENTAL_REFRESH
	var respMsg *quickfix.Message

	go func() {
		time.Sleep(responseTime)
		// mock result for marketDataRequest response
		msgStr := "8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Haha\u0001268=0\u000110=196\u0001"
		respMsg = getMsgFromString(msgStr)

		mutex.Lock()
		respMsg.Body.Set(field.NewMDReqID(requestID))
		mutex.Unlock()

		err := mockDeribitResponse(respMsg)
		require.NoError(err)
	}()

	msg, err := ts.c.MarketDataRequest(
		context.Background(),
		enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES,
		&marketDepth,
		&mdUpdateType,
		[]enum.MDEntryType{
			enum.MDEntryType_BID,
			enum.MDEntryType_OFFER,
		},
		[]string{"BTC-26AUG22-29000-C, BTC-PERPETUAL"},
	)
	require.NoError(err)
	require.Equal(msg.String(), respMsg.String())
}

// nolint:lll
func (ts *FixTestSuite) TestSubscribeOrderBooks() {
	require := ts.Require()
	instruments := []string{"BTC-26AUG22-29000-C, BTC-PERPETUAL"}
	wrongInstruments := []string{"SHIB-26AUG22-0.000123-C, SHIB-PERPETUAL"}
	sendResp := make(chan bool)

	go func() {
		msgStrings := []string{
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
			"8=FIX.4.4\u00019=189\u000135=Y\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
		}
		for _, msgStr := range msgStrings {
			<-sendResp
			// subscribe orderbook mock response
			time.Sleep(responseTime)
			respMsg := getMsgFromString(msgStr)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}
	}()

	sendResp <- true
	err := ts.c.SubscribeOrderBooks(context.Background(), instruments)
	require.NoError(err)

	sendResp <- true
	err = ts.c.SubscribeOrderBooks(context.Background(), wrongInstruments)
	require.Error(err) // MsgType_MARKET_DATA_REQUEST_REJECT
}

// nolint:lll
func (ts *FixTestSuite) TestUnsubscribeOrderBooks() {
	require := ts.Require()
	instruments := []string{"BTC-26AUG22-29000-C, BTC-PERPETUAL"}
	wrongInstruments := []string{"SHIB-26AUG22-0.000123-C, SHIB-PERPETUAL"}
	sendResp := make(chan bool)

	go func() {
		msgStrings := []string{
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
			"8=FIX.4.4\u00019=189\u000135=Y\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
		}
		for _, msgStr := range msgStrings {
			<-sendResp
			// unsubscribe orderbook mock response
			time.Sleep(responseTime)
			respMsg := getMsgFromString(msgStr)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}
	}()

	sendResp <- true
	err := ts.c.UnsubscribeOrderBooks(context.Background(), instruments)
	require.NoError(err)

	sendResp <- true
	err = ts.c.UnsubscribeOrderBooks(context.Background(), wrongInstruments)
	require.Error(err) // MsgType_MARKET_DATA_REQUEST_REJECT
}

// nolint:lll
func (ts *FixTestSuite) TestSubscribeTrades() {
	require := ts.Require()
	instruments := []string{"BTC-26AUG22-29000-C, BTC-PERPETUAL"}
	wrongInstruments := []string{"SHIB-26AUG22-0.000123-C, SHIB-PERPETUAL"}
	sendResp := make(chan bool)

	go func() {
		msgStrings := []string{
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
			"8=FIX.4.4\u00019=189\u000135=Y\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
		}
		for _, msgStr := range msgStrings {
			<-sendResp
			// subscribe trades mock response
			time.Sleep(responseTime)
			respMsg := getMsgFromString(msgStr)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}
	}()

	sendResp <- true
	err := ts.c.SubscribeTrades(context.Background(), instruments)
	require.NoError(err)

	sendResp <- true
	err = ts.c.SubscribeTrades(context.Background(), wrongInstruments)
	require.Error(err) // MsgType_MARKET_DATA_REQUEST_REJECT
}

// nolint:lll
func (ts *FixTestSuite) TestUnsubscribeTrades() {
	require := ts.Require()
	instruments := []string{"BTC-26AUG22-29000-C, BTC-PERPETUAL"}
	wrongInstruments := []string{"SHIB-26AUG22-0.000123-C, SHIB-PERPETUAL"}
	sendResp := make(chan bool)

	go func() {
		msgStrings := []string{
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
			"8=FIX.4.4\u00019=189\u000135=Y\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hoho\u0001268=0\u000110=196\u0001",
		}
		for _, msgStr := range msgStrings {
			<-sendResp
			// unsubscribe trades mock response
			time.Sleep(responseTime)
			respMsg := getMsgFromString(msgStr)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}
	}()

	sendResp <- true
	err := ts.c.UnsubscribeTrades(context.Background(), instruments)
	require.NoError(err)

	sendResp <- true
	err = ts.c.UnsubscribeTrades(context.Background(), wrongInstruments)
	require.Error(err) // MsgType_MARKET_DATA_REQUEST_REJECT
}

// nolint:lll,dupl
func (ts *FixTestSuite) TestSubscribe() {
	require := ts.Require()

	type testSubscribe struct {
		channel       string
		fixResp       string
		expectedError error
	}

	tests := []testSubscribe{
		{
			"book.BTC-PERPETUAL",
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hiha\u0001268=0\u000110=196\u0001",
			nil,
		},
		{
			"trades.BTC-PERPETUAL",
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hehe\u0001268=0\u000110=196\u0001",
			nil,
		},
	}

	require.Len(ts.c.subscriptions, 0)

	for _, test := range tests {
		// subscribe mock response
		go func(test testSubscribe) {
			time.Sleep(responseTime)
			respMsg := getMsgFromString(test.fixResp)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}(test)

		err := ts.c.Subscribe(context.Background(), []string{test.channel})
		require.ErrorIs(err, test.expectedError)
	}

	require.Len(ts.c.subscriptions, 2)
}

// nolint:lll,dupl
func (ts *FixTestSuite) TestUnsubscribe() {
	require := ts.Require()

	type testUnsubscribe struct {
		channel       string
		fixResp       string
		expectedError error
	}

	tests := []testUnsubscribe{
		{
			"book.BTC-PERPETUAL",
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Hihu\u0001268=0\u000110=196\u0001",
			nil,
		},
		{
			"trades.BTC-PERPETUAL",
			"8=FIX.4.4\u00019=189\u000135=W\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=2\u000152=20220822-03:44:38.119\u000155=BTC-26AUG22-22500-P\u0001231=1.0000\u0001311=BTC-26AUG22\u0001810=21443.3568\u0001100087=0.0000\u0001100090=0.0504\u0001746=0.0000\u0001201=0\u0001262=Huha\u0001268=0\u000110=196\u0001",
			nil,
		},
	}

	require.Len(ts.c.subscriptions, 2)

	for _, test := range tests {
		// unsubscribe mock response
		go func(test testUnsubscribe) {
			time.Sleep(responseTime)
			respMsg := getMsgFromString(test.fixResp)

			mutex.Lock()
			respMsg.Body.Set(field.NewMDReqID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}(test)

		err := ts.c.Unsubscribe(context.Background(), []string{test.channel})
		require.ErrorIs(err, test.expectedError)
	}

	require.Len(ts.c.subscriptions, 0)
}

// nolint:lll,funlen
func (ts *FixTestSuite) TestCreateOrder() {
	require := ts.Require()
	sendResp := make(chan bool)

	go func() {
		msgStrings := []string{
			"8=FIX.4.4\u00019=494\u000135=8\u000149=DERIBITSERVER\u000156=FIX_TEST\u000134=9158\u000152=20220818-06:37:42.584\u0001527=14020845373\u000137=14020845373\u000111=14020845373\u000141=6b5c1fe2-e6ad-4ccf-93d7-8b6ccd51cdea\u0001150=I\u000139=2\u000154=1\u000160=20220818-06:37:42.583\u000112=0.00003000\u0001151=0.0000\u000114=0.1000\u000138=0.1000\u000140=2\u000144=0.077\u0001103=0\u000158=success\u0001207=DERIBITSERVER\u000155=BTC-19AUG22-21000-C\u0001854=1\u0001231=1.0000\u00016=0.077000\u0001210=0.1000\u0001100010=BTC-19AUG22-21000-C_buy_0.077_0.1_JNGhwLYqkJzoYSk\u000132=0.1000\u000131=0.0770\u00011362=1\u00011363=BTC-19AUG22-21000-C#102\u00011364=0.0770\u00011365=0.1000\u00011443=2\u000110=012\u0001",
			"8=FIX.4.4\u00019=432\u000135=8\u000149=DERIBITSERVER\u000156=OPTION_TRADING_TEST\u000134=7534\u000152=20220912-11:12:19.623\u0001527=14230452591\u000137=14230452591\u000111=14230452591\u000141=ff39b93f-4e8c-4187-8930-6eac35bdb2e9\u0001150=I\u000139=2\u000154=1\u000160=20220912-11:12:19.623\u000112=0.00300000\u0001151=0.0\u000114=10.0\u000138=10.0\u000140=2\u000144=0.0865\u0001103=0\u000158=success\u0001207=DERIBITSERVER\u000155=BTC-30JUN23-14000-P\u0001854=1\u0001231=1.0\u00016=0.086500\u0001210=10.0\u000132=10.0\u000131=0.0865\u00011362=1\u00011363=BTC-30JUN23-14000-P#83\u00011364=0.0865\u00011365=10.0\u00011443=2\u000110=071\u0001",
		}
		for _, msgStr := range msgStrings {
			<-sendResp
			// create order mock response
			time.Sleep(responseTime)
			respMsg := getMsgFromString(msgStr)

			mutex.Lock()
			respMsg.Body.Set(field.NewOrigClOrdID(requestID))
			mutex.Unlock()

			err := mockDeribitResponse(respMsg)
			require.NoError(err)
		}
	}()

	// success case
	sendResp <- true
	res, err := ts.c.CreateOrder(
		context.Background(),
		"BTC-19AUG22-21000-C",             // symbol
		"buy",                             // side
		0.1,                               // amount
		0.077,                             // price
		enum.OrdType_LIMIT,                // order_type
		enum.TimeInForce_GOOD_TILL_CANCEL, // time_in_force
		"",                                // execInst
		"BTC-19AUG22-21000-C_buy_0.077_0.1_JNGhwLYqkJzoYSk", // CltOrdId
	)
	expectedOutput := models.Order{
		OrderState:          "filled",
		MaxShow:             0.1,
		API:                 true,
		Amount:              0.1,
		Web:                 false,
		InstrumentName:      "BTC-19AUG22-21000-C",
		OriginalOrderType:   "limit",
		Price:               0.077,
		TimeInForce:         "good_til_cancelled",
		LastUpdateTimestamp: 1660804662583,
		PostOnly:            false,
		Replaced:            false,
		FilledAmount:        0.1,
		AveragePrice:        0.077,
		OrderID:             "14020845373",
		ReduceOnly:          false,
		Commission:          3e-05,
		Label:               "BTC-19AUG22-21000-C_buy_0.077_0.1_JNGhwLYqkJzoYSk",
		CreationTimestamp:   1660804662583,
		Direction:           "buy",
		OrderType:           "limit",
	}

	require.NoError(err)
	require.Equal(res, expectedOutput)

	// error case
	sendResp <- true
	res, err = ts.c.CreateOrder(
		context.Background(),
		"BTC-19AUG22-21000-C",             // symbol
		"buy",                             // side
		0.1,                               // amount
		0.077,                             // price
		enum.OrdType_LIMIT,                // order_type
		enum.TimeInForce_GOOD_TILL_CANCEL, // time_in_force
		"",                                // execInst
		"BTC-19AUG22-21000-C_buy_0.077_0.1_JNGhwLYqkJzoYSk", // CltOrdId
	)
	require.Error(err)
	require.Equal(res, models.Order{})
}

func getMsgFromString(str string) *quickfix.Message {
	msg := quickfix.NewMessage()
	bufferData := bytes.NewBufferString(str)
	err := quickfix.ParseMessage(msg, bufferData)
	if err != nil {
		return quickfix.NewMessage()
	}
	return msg
}

func (ts *FixTestSuite) TestXClose() {
	require := ts.Require()

	// subscribe book and trades channels
	require.Len(ts.c.subscriptions, 0)
	ts.TestSubscribe()
	require.Len(ts.c.subscriptions, 2)
	require.Len(ts.c.subscriptionsMap, 2)

	require.True(ts.c.IsConnected())
	ts.c.Close()
	require.False(ts.c.IsConnected())
}
