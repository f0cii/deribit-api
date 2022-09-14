package websocket

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/getlantern/deepcopy"
	"github.com/sourcegraph/jsonrpc2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	successResponse = "success"
)

type MockRPCConn struct {
	addr         string
	handler      jsonrpc2.Handler
	disconnectCh chan struct{}
	results      []interface{}
}

func NewMockRCConn(ctx context.Context, addr string, h jsonrpc2.Handler) (JSONRPC2, error) {
	return &MockRPCConn{
		addr:         addr,
		handler:      h,
		disconnectCh: make(chan struct{}),
		results: []interface{}{
			&models.AuthResponse{},
			"success",
		},
	}, nil
}

func (c *MockRPCConn) Call(
	ctx context.Context,
	method string,
	params interface{},
	result interface{},
	opt ...jsonrpc2.CallOption,
) error {
	res := c.GetResult()
	if res == nil {
		return nil
	}

	return deepcopy.Copy(result, res)
}

func (c *MockRPCConn) Notify(
	ctx context.Context,
	method string,
	params interface{},
	opt ...jsonrpc2.CallOption,
) error {
	return nil
}

func (c *MockRPCConn) Close() error {
	close(c.disconnectCh)
	return nil
}

func (c *MockRPCConn) DisconnectNotify() <-chan struct{} {
	return c.disconnectCh
}

func (c *MockRPCConn) AddResult(result interface{}) {
	c.results = append(c.results, result)
}

func (c *MockRPCConn) GetResult() interface{} {
	if len(c.results) == 0 {
		return nil
	}

	res := c.results[0]
	c.results = c.results[1:]
	return res
}

func addResult(conn JSONRPC2, res interface{}) {
	conn.(*MockRPCConn).AddResult(res)
}

func newClient() *Client {
	cfg := Configuration{
		Addr:          TestBaseURL,
		APIKey:        "test_api_key",
		SecretKey:     "test_secret_key",
		DebugMode:     true,
		NewRPCConn:    NewMockRCConn,
		AutoReconnect: true,
	}

	return New(zap.S(), &cfg)
}

// nolint:gochecknoglobals
var testClient *Client

func TestMain(m *testing.M) {
	testClient = newClient()
	if err := testClient.Start(); err != nil {
		panic(err)
	}

	retCode := m.Run()

	testClient.Stop()

	os.Exit(retCode)
}

func TestNewRPCConn(t *testing.T) {
	conn, err := NewRPCConn(context.Background(), TestBaseURL, nil)
	if assert.NoError(t, err) {
		assert.NotNil(t, conn)
	}
}

func TestNewClient(t *testing.T) {
	cfg := Configuration{
		Addr:       TestBaseURL,
		DebugMode:  true,
		NewRPCConn: nil,
	}
	client := New(zap.S(), &cfg)
	assert.NotNil(t, client)
}

func TestStartStop(t *testing.T) {
	client := newClient()
	err := client.Start()
	require.NoError(t, err)
	assert.True(t, client.IsConnected())
	client.Stop()
	assert.False(t, client.IsConnected())
}

func TestCall(t *testing.T) {
	client := newClient()
	err := client.Call(context.Background(), "public/test", nil, nil)
	assert.ErrorIs(t, ErrNotConnected, err)

	err = client.Start()
	require.NoError(t, err)
	addResult(client.rpcConn, &models.TestResponse{Version: "1.2.26"})

	var testResp models.TestResponse
	err = client.Call(context.Background(), "public/test", nil, &testResp)
	if assert.NoError(t, err) {
		assert.Equal(t, "1.2.26", testResp.Version)
	}
}

// nolint:lll,funlen,maintidx
func TestHandle(t *testing.T) {
	tests := []struct {
		req    *jsonrpc2.Request
		params Event
		expect interface{}
	}{
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "announcements",
				Data:    nil,
			},
			expect: &models.AnnouncementsNotification{},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "book.BTC-PERPETUAL.raw",
				Data:    json.RawMessage("{\"timestamp\":1662714568585,\"prev_change_id\":14214947552,\"instrument_name\":\"BTC-PERPETUAL\",\"change_id\":14214947618,\"bids\":[[\"new\",20338,20700],[\"delete\",20337,0]],\"asks\":[[\"change\",20644.5,2580],[\"new\",20684,3510],[\"delete\",20686.5,0]]}"),
			},
			expect: &models.OrderBookRawNotification{
				Timestamp:      1662714568585,
				InstrumentName: "BTC-PERPETUAL",
				PrevChangeID:   14214947552,
				ChangeID:       14214947618,
				Bids: []models.OrderBookNotificationItem{
					{
						Action: "new",
						Price:  20338,
						Amount: 20700,
					},
					{
						Action: "delete",
						Price:  20337,
					},
				},
				Asks: []models.OrderBookNotificationItem{
					{
						Action: "change",
						Price:  20644.5,
						Amount: 2580,
					},
					{
						Action: "new",
						Price:  20684,
						Amount: 3510,
					},
					{
						Action: "delete",
						Price:  20686.5,
					},
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "book.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("{\"type\":\"change\",\"timestamp\":1662714568585,\"prev_change_id\":14214947552,\"instrument_name\":\"BTC-PERPETUAL\",\"change_id\":14214947618,\"bids\":[[\"new\",20338,20700],[\"delete\",20337,0]],\"asks\":[[\"change\",20644.5,2580],[\"new\",20684,3510],[\"delete\",20686.5,0]]}"),
			},
			expect: &models.OrderBookNotification{
				Type:           "change",
				Timestamp:      1662714568585,
				InstrumentName: "BTC-PERPETUAL",
				PrevChangeID:   14214947552,
				ChangeID:       14214947618,
				Bids: []models.OrderBookNotificationItem{
					{
						Action: "new",
						Price:  20338,
						Amount: 20700,
					},
					{
						Action: "delete",
						Price:  20337,
					},
				},
				Asks: []models.OrderBookNotificationItem{
					{
						Action: "change",
						Price:  20644.5,
						Amount: 2580,
					},
					{
						Action: "new",
						Price:  20684,
						Amount: 3510,
					},
					{
						Action: "delete",
						Price:  20686.5,
					},
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "book.BTC-PERPETUAL.none.1.100ms",
				Data:    json.RawMessage("{\"timestamp\":1662715579344,\"instrument_name\":\"BTC-PERPETUAL\",\"change_id\":14214997020,\"bids\":[[20659,3970]],\"asks\":[[20661.5,190]]}"),
			},
			expect: &models.OrderBookGroupNotification{
				Timestamp:      1662715579344,
				InstrumentName: "BTC-PERPETUAL",
				ChangeID:       14214997020,
				Bids:           [][]float64{{20659, 3970}},
				Asks:           [][]float64{{20661.5, 190}},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "deribit_price_index.btc_usd",
				Data:    json.RawMessage("{\"timestamp\":1662715972131,\"price\":20651.5,\"index_name\":\"btc_usd\"}"),
			},
			expect: &models.DeribitPriceIndexNotification{
				Timestamp: 1662715972131,
				Price:     20651.5,
				IndexName: "btc_usd",
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "deribit_price_ranking.btc_usd",
				Data:    json.RawMessage("[{\"weight\":0,\"timestamp\":1662716219588,\"price\":20674.5,\"original_price\":20674.5,\"identifier\":\"bitfinex\",\"enabled\":false},{\"weight\":12.5,\"timestamp\":1662716219734,\"price\":20669.5,\"original_price\":20669.5,\"identifier\":\"bitstamp\",\"enabled\":true},{\"weight\":0,\"timestamp\":1662716216232,\"price\":20668.85,\"original_price\":20668.85,\"identifier\":\"bittrex\",\"enabled\":false},{\"weight\":12.5,\"timestamp\":1662716219814,\"price\":20669.89,\"original_price\":20669.89,\"identifier\":\"coinbase\",\"enabled\":true},{\"weight\":12.5,\"timestamp\":1662716219713,\"price\":20668.5,\"original_price\":20668.5,\"identifier\":\"ftx\",\"enabled\":true},{\"weight\":12.5,\"timestamp\":1662716219375,\"price\":20669.42,\"original_price\":20669.42,\"identifier\":\"gateio\",\"enabled\":true},{\"weight\":12.5,\"timestamp\":1662716219782,\"price\":20668.4,\"original_price\":20668.4,\"identifier\":\"gemini\",\"enabled\":true},{\"weight\":12.5,\"timestamp\":1662716218446,\"price\":20665.38,\"original_price\":20665.38,\"identifier\":\"itbit\",\"enabled\":true},{\"weight\":12.5,\"timestamp\":1662716217419,\"price\":20655.75,\"original_price\":20655.75,\"identifier\":\"kraken\",\"enabled\":true},{\"weight\":0,\"timestamp\":1662716219562,\"price\":20674,\"original_price\":20674,\"identifier\":\"lmax\",\"enabled\":false},{\"weight\":12.5,\"timestamp\":1662716219000,\"price\":20666.29,\"original_price\":20666.29,\"identifier\":\"okcoin\",\"enabled\":true}]"),
			},
			expect: &models.DeribitPriceRankingNotification{
				{
					Weight:     0,
					Timestamp:  1662716219588,
					Price:      20674.5,
					Identifier: "bitfinex",
					Enabled:    false,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219734,
					Price:      20669.5,
					Identifier: "bitstamp",
					Enabled:    true,
				},
				{
					Weight:     0,
					Timestamp:  1662716216232,
					Price:      20668.85,
					Identifier: "bittrex",
					Enabled:    false,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219814,
					Price:      20669.89,
					Identifier: "coinbase",
					Enabled:    true,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219713,
					Price:      20668.5,
					Identifier: "ftx",
					Enabled:    true,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219375,
					Price:      20669.42,
					Identifier: "gateio",
					Enabled:    true,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219782,
					Price:      20668.4,
					Identifier: "gemini",
					Enabled:    true,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716218446,
					Price:      20665.38,
					Identifier: "itbit",
					Enabled:    true,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716217419,
					Price:      20655.75,
					Identifier: "kraken",
					Enabled:    true,
				},
				{
					Weight:     0,
					Timestamp:  1662716219562,
					Price:      20674,
					Identifier: "lmax",
					Enabled:    false,
				},
				{
					Weight:     12.5,
					Timestamp:  1662716219000,
					Price:      20666.29,
					Identifier: "okcoin",
					Enabled:    true,
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "estimated_expiration_price.btc_usd",
				Data:    json.RawMessage("{\"seconds\":76228,\"price\":21094.14,\"is_estimated\":false}"),
			},
			expect: &models.EstimatedExpirationPriceNotification{
				Seconds:     76228,
				Price:       21094.14,
				IsEstimated: false,
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "markprice.options.btc_usd",
				Data:    json.RawMessage("[{\"timestamp\":1662720695772,\"mark_price\":0.1523,\"iv\":0.7686,\"instrument_name\":\"BTC-31MAR23-26000-C\"},{\"timestamp\":1662720695772,\"mark_price\":0.1542,\"iv\":0.6748,\"instrument_name\":\"BTC-28OCT22-23000-P\"}]"),
			},
			expect: &models.MarkpriceOptionsNotification{
				{
					MarkPrice:      0.1523,
					Iv:             0.7686,
					InstrumentName: "BTC-31MAR23-26000-C",
				},
				{
					MarkPrice:      0.1542,
					Iv:             0.6748,
					InstrumentName: "BTC-28OCT22-23000-P",
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "perpetual.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("{\"timestamp\":1662721149878,\"interest\":-0.005,\"index_price\":21073.99}"),
			},
			expect: &models.PerpetualNotification{
				Timestamp:  1662721149878,
				Interest:   -0.005,
				IndexPrice: 21073.99,
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "quote.BTC-PERPETUAL",
				Data:    json.RawMessage("{\"timestamp\":1662721273742,\"instrument_name\":\"BTC-PERPETUAL\",\"best_bid_price\":21070,\"best_bid_amount\":1010,\"best_ask_price\":21075,\"best_ask_amount\":3730}"),
			},
			expect: &models.QuoteNotification{
				Timestamp:      1662721273742,
				InstrumentName: "BTC-PERPETUAL",
				BestBidPrice:   float64Pointer(21070),
				BestBidAmount:  1010,
				BestAskPrice:   float64Pointer(21075),
				BestAskAmount:  3730,
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "ticker.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("{\"timestamp\":1662721394017,\"stats\":{\"volume_usd\":194393720,\"volume\":9678,\"price_change\":8.5,\"low\":19025,\"high\":21100.5},\"state\":\"open\",\"settlement_price\":20585.66,\"open_interest\":2970900920,\"min_price\":20685.71,\"max_price\":21315.73,\"mark_price\":20992.53,\"last_price\":21007.5,\"interest_value\":-1.5464265205734715,\"instrument_name\":\"BTC-PERPETUAL\",\"index_price\":21025.44,\"funding_8h\":-0.001,\"estimated_delivery_price\":21025.44,\"current_funding\":-0.001,\"best_bid_price\":20986,\"best_bid_amount\":460,\"best_ask_price\":20987.5,\"best_ask_amount\":400}"),
			},
			expect: &models.TickerNotification{
				Timestamp: 1662721394017,
				Stats: models.Stats{
					Volume:      9678,
					PriceChange: float64Pointer(8.5),
					Low:         19025,
					High:        21100.5,
				},
				State:           "open",
				SettlementPrice: 20585.66,
				OpenInterest:    2970900920,
				MinPrice:        20685.71,
				MaxPrice:        21315.73,
				MarkPrice:       20992.53,
				LastPrice:       21007.5,
				InstrumentName:  "BTC-PERPETUAL",
				IndexPrice:      21025.44,
				Funding8H:       -0.001,
				CurrentFunding:  -0.001,
				BestBidPrice:    float64Pointer(20986),
				BestBidAmount:   460,
				BestAskPrice:    float64Pointer(20987.5),
				BestAskAmount:   400,
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "trades.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("[ { \"trade_seq\": 81769518, \"trade_id\": \"119810484\", \"timestamp\": 1662957035112, \"tick_direction\": 2, \"price\": 21703, \"mark_price\": 21705.36, \"instrument_name\": \"BTC-PERPETUAL\", \"index_price\": 21711.76, \"direction\": \"sell\", \"amount\": 1000 } ]"),
			},
			expect: &models.TradesNotification{
				{
					Amount:         1000,
					BlockTradeID:   "",
					Direction:      "sell",
					IndexPrice:     21711.76,
					InstrumentName: "BTC-PERPETUAL",
					MarkPrice:      21705.36,
					Price:          21703,
					TickDirection:  2,
					Timestamp:      1662957035112,
					TradeID:        "119810484",
					TradeSeq:       81769518,
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "user.changes.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("{ \"trades\": [ { \"trade_seq\": 81772419, \"trade_id\": \"119813642\", \"timestamp\": 1662964399064, \"tick_direction\": 0, \"state\": \"filled\", \"self_trade\": false, \"risk_reducing\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 21760.5, \"post_only\": false, \"order_type\": \"market\", \"order_id\": \"14228823973\", \"mmp\": false, \"matching_id\": null, \"mark_price\": 21758.49, \"liquidity\": \"T\", \"instrument_name\": \"BTC-PERPETUAL\", \"index_price\": 21755.39, \"fee_currency\": \"BTC\", \"fee\": 4.6e-7, \"direction\": \"buy\", \"api\": false, \"amount\": 20 }, { \"trade_seq\": 81772420, \"trade_id\": \"119813643\", \"timestamp\": 1662964399064, \"tick_direction\": 1, \"state\": \"filled\", \"self_trade\": false, \"risk_reducing\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 21760.5, \"post_only\": false, \"order_type\": \"market\", \"order_id\": \"14228823973\", \"mmp\": false, \"matching_id\": null, \"mark_price\": 21758.49, \"liquidity\": \"T\", \"instrument_name\": \"BTC-PERPETUAL\", \"index_price\": 21755.39, \"fee_currency\": \"BTC\", \"fee\": 0.00000184, \"direction\": \"buy\", \"api\": false, \"amount\": 80 } ], \"positions\": [ { \"total_profit_loss\": -4.24e-7, \"size_currency\": 0.004595907, \"size\": 100, \"settlement_price\": 21623.41, \"realized_profit_loss\": 0, \"realized_funding\": 0, \"open_orders_margin\": 0, \"mark_price\": 21758.49, \"maintenance_margin\": 0.00004596, \"leverage\": 50, \"kind\": \"future\", \"interest_value\": -9.458785481892445, \"instrument_name\": \"BTC-PERPETUAL\", \"initial_margin\": 0.000091919, \"index_price\": 21755.39, \"floating_profit_loss\": -4.24e-7, \"direction\": \"buy\", \"delta\": 0.004595907, \"average_price\": 21760.5 } ], \"orders\": [ { \"web\": true, \"time_in_force\": \"good_til_cancelled\", \"risk_reducing\": false, \"replaced\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 22084, \"post_only\": false, \"order_type\": \"market\", \"order_state\": \"filled\", \"order_id\": \"14228823973\", \"mmp\": false, \"max_show\": 100, \"last_update_timestamp\": 1662964399064, \"label\": \"\", \"is_liquidation\": false, \"instrument_name\": \"BTC-PERPETUAL\", \"filled_amount\": 100, \"direction\": \"buy\", \"creation_timestamp\": 1662964399064, \"commission\": 0.0000023, \"average_price\": 21760.5, \"api\": false, \"amount\": 100 } ], \"instrument_name\": \"BTC-PERPETUAL\" }"),
			},
			expect: &models.UserChangesNotification{
				InstrumentName: "BTC-PERPETUAL",
				Trades: []models.UserTrade{
					{
						UnderlyingPrice: 0,
						TradeSeq:        0x4dfbf83,
						TradeID:         "119813642",
						Timestamp:       0x1833066fbd8,
						TickDirection:   0,
						State:           "filled",
						SelfTrade:       false,
						ReduceOnly:      false,
						ProfitLost:      0,
						Price:           21760.5,
						PostOnly:        false,
						OrderType:       "market",
						OrderID:         "14228823973",
						MatchingID:      nil,
						MarkPrice:       21758.49,
						Liquidity:       "T",
						Liquidation:     "",
						Label:           "",
						IV:              0,
						InstrumentName:  "BTC-PERPETUAL",
						IndexPrice:      21755.39,
						FeeCurrency:     "BTC",
						Fee:             4.6e-07,
						Direction:       "buy",
						Amount:          20,
						BlockTradeID:    "",
					},
					{
						UnderlyingPrice: 0,
						TradeSeq:        0x4dfbf84,
						TradeID:         "119813643",
						Timestamp:       0x1833066fbd8,
						TickDirection:   1,
						State:           "filled",
						SelfTrade:       false,
						ReduceOnly:      false,
						ProfitLost:      0,
						Price:           21760.5,
						PostOnly:        false,
						OrderType:       "market",
						OrderID:         "14228823973",
						MatchingID:      nil,
						MarkPrice:       21758.49,
						Liquidity:       "T",
						Liquidation:     "",
						Label:           "",
						IV:              0,
						InstrumentName:  "BTC-PERPETUAL",
						IndexPrice:      21755.39,
						FeeCurrency:     "BTC",
						Fee:             1.84e-06,
						Direction:       "buy",
						Amount:          80,
						BlockTradeID:    "",
					},
				},
				Positions: []models.Position{
					{
						AveragePrice:              21760.5,
						AveragePriceUSD:           0,
						Delta:                     0.004595907,
						Direction:                 "buy",
						EstimatedLiquidationPrice: 0,
						FloatingProfitLoss:        -4.24e-07,
						FloatingProfitLossUSD:     0,
						Gamma:                     0,
						IndexPrice:                21755.39,
						InitialMargin:             9.1919e-05,
						InstrumentName:            "BTC-PERPETUAL",
						Kind:                      "future",
						Leverage:                  50,
						MaintenanceMargin:         4.596e-05,
						MarkPrice:                 21758.49,
						OpenOrdersMargin:          0,
						RealizedFunding:           0,
						RealizedProfitLoss:        0,
						SettlementPrice:           21623.41,
						Size:                      100,
						SizeCurrency:              0.004595907,
						Theta:                     0,
						TotalProfitLoss:           -4.24e-07,
						Vega:                      0,
					},
				},
				Orders: []models.Order{
					{
						MMPCancelled:        false,
						OrderState:          "filled",
						MaxShow:             100,
						API:                 false,
						Amount:              100,
						Web:                 true,
						InstrumentName:      "BTC-PERPETUAL",
						Advanced:            "",
						Triggered:           nil,
						BlockTrade:          false,
						OriginalOrderType:   "",
						Price:               22084,
						TimeInForce:         "good_til_cancelled",
						AutoReplaced:        false,
						LastUpdateTimestamp: 0x1833066fbd8,
						PostOnly:            false,
						Replaced:            false,
						FilledAmount:        100,
						AveragePrice:        21760.5,
						OrderID:             "14228823973",
						ReduceOnly:          false,
						Commission:          2.3e-06,
						AppName:             "",
						Label:               "",
						TriggerOrderID:      "",
						TriggedPrice:        0,
						CreationTimestamp:   0x1833066fbd8,
						Direction:           "buy",
						IsLiquidation:       false,
						OrderType:           "market",
						USD:                 0,
						ProfitLoss:          0,
						Implv:               0,
						Trigger:             "",
					},
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "user.orders.BTC-PERPETUAL.raw",
				Data:    json.RawMessage("{ \"web\": true, \"time_in_force\": \"good_til_cancelled\", \"risk_reducing\": false, \"replaced\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 22084, \"post_only\": false, \"order_type\": \"market\", \"order_state\": \"filled\", \"order_id\": \"14228823973\", \"mmp\": false, \"max_show\": 100, \"last_update_timestamp\": 1662964399064, \"label\": \"\", \"is_liquidation\": false, \"instrument_name\": \"BTC-PERPETUAL\", \"filled_amount\": 100, \"direction\": \"buy\", \"creation_timestamp\": 1662964399064, \"commission\": 0.0000023, \"average_price\": 21760.5, \"api\": false, \"amount\": 100 }"),
			},
			expect: &models.Order{
				OrderState:          "filled",
				MaxShow:             100,
				Amount:              100,
				Web:                 true,
				InstrumentName:      "BTC-PERPETUAL",
				Price:               22084,
				TimeInForce:         "good_til_cancelled",
				LastUpdateTimestamp: 1662964399064,
				FilledAmount:        100,
				AveragePrice:        21760.5,
				OrderID:             "14228823973",
				Commission:          0.0000023,
				CreationTimestamp:   1662964399064,
				Direction:           "buy",
				OrderType:           "market",
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "user.orders.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("[ { \"web\": true, \"time_in_force\": \"good_til_cancelled\", \"risk_reducing\": false, \"replaced\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 22084, \"post_only\": false, \"order_type\": \"market\", \"order_state\": \"filled\", \"order_id\": \"14228823973\", \"mmp\": false, \"max_show\": 100, \"last_update_timestamp\": 1662964399064, \"label\": \"\", \"is_liquidation\": false, \"instrument_name\": \"BTC-PERPETUAL\", \"filled_amount\": 100, \"direction\": \"buy\", \"creation_timestamp\": 1662964399064, \"commission\": 0.0000023, \"average_price\": 21760.5, \"api\": false, \"amount\": 100 } ]"),
			},
			expect: &models.UserOrderNotification{
				{
					OrderState:          "filled",
					MaxShow:             100,
					Amount:              100,
					Web:                 true,
					InstrumentName:      "BTC-PERPETUAL",
					Price:               22084,
					TimeInForce:         "good_til_cancelled",
					LastUpdateTimestamp: 1662964399064,
					FilledAmount:        100,
					AveragePrice:        21760.5,
					OrderID:             "14228823973",
					Commission:          0.0000023,
					CreationTimestamp:   1662964399064,
					Direction:           "buy",
					OrderType:           "market",
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "user.portfolio.BTC",
				Data:    json.RawMessage("{ \"total_pl\": -107.54329243, \"session_upl\": 2.03600535, \"session_rpl\": 0, \"projected_maintenance_margin\": 104.88127092, \"projected_initial_margin\": 115.13608906, \"projected_delta_total\": 288.6408, \"portfolio_margining_enabled\": false, \"options_vega\": 2697.2679, \"options_value\": 6.43407021, \"options_theta\": -2177.2495, \"options_session_upl\": 0.54354297, \"options_session_rpl\": 0, \"options_pl\": 27.59762833, \"options_gamma\": 0.0043, \"options_delta\": -5.4052, \"margin_balance\": 6663.85236718, \"maintenance_margin\": 104.88127092, \"initial_margin\": 115.13608906, \"futures_session_upl\": 1.49246237, \"futures_session_rpl\": 0, \"futures_pl\": -135.14092076, \"fee_balance\": 0, \"estimated_liquidation_ratio_map\": { \"btc_usd\": 0.05449954798610861 }, \"estimated_liquidation_ratio\": 0.05449955, \"equity\": 6670.28643738, \"delta_total_map\": { \"btc_usd\": 300.48007431300005 }, \"delta_total\": 288.6408, \"currency\": \"BTC\", \"balance\": 6662.3599048, \"available_withdrawal_funds\": 6547.22381574, \"available_funds\": 6548.71627812 }"),
			},
			expect: &models.PortfolioNotification{
				TotalPl:                    -107.54329243,
				SessionUpl:                 2.03600535,
				SessionRpl:                 0,
				ProjectedMaintenanceMargin: 104.88127092,
				ProjectedInitialMargin:     115.13608906,
				ProjectedDeltaTotal:        288.6408,
				PortfolioMarginingEnabled:  false,
				OptionsVega:                2697.2679,
				OptionsValue:               6.43407021,
				OptionsTheta:               -2177.2495,
				OptionsSessionUpl:          0.54354297,
				OptionsSessionRpl:          0,
				OptionsPl:                  27.59762833,
				OptionsGamma:               0.0043,
				OptionsDelta:               -5.4052,
				MarginBalance:              6663.85236718,
				MaintenanceMargin:          104.88127092,
				InitialMargin:              115.13608906,
				FuturesSessionUpl:          1.49246237,
				FuturesSessionRpl:          0,
				FuturesPl:                  -135.14092076,
				EstimatedLiquidationRatio:  0.05449955,
				Equity:                     6670.28643738,
				DeltaTotal:                 288.6408,
				Currency:                   "BTC",
				Balance:                    6662.3599048,
				AvailableWithdrawalFunds:   6547.22381574,
				AvailableFunds:             6548.71627812,
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "user.trades.BTC-PERPETUAL.100ms",
				Data:    json.RawMessage("[ { \"trade_seq\": 81772419, \"trade_id\": \"119813642\", \"timestamp\": 1662964399064, \"tick_direction\": 0, \"state\": \"filled\", \"self_trade\": false, \"risk_reducing\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 21760.5, \"post_only\": false, \"order_type\": \"market\", \"order_id\": \"14228823973\", \"mmp\": false, \"matching_id\": null, \"mark_price\": 21758.49, \"liquidity\": \"T\", \"instrument_name\": \"BTC-PERPETUAL\", \"index_price\": 21755.39, \"fee_currency\": \"BTC\", \"fee\": 4.6e-7, \"direction\": \"buy\", \"api\": false, \"amount\": 20 }, { \"trade_seq\": 81772420, \"trade_id\": \"119813643\", \"timestamp\": 1662964399064, \"tick_direction\": 1, \"state\": \"filled\", \"self_trade\": false, \"risk_reducing\": false, \"reduce_only\": false, \"profit_loss\": 0, \"price\": 21760.5, \"post_only\": false, \"order_type\": \"market\", \"order_id\": \"14228823973\", \"mmp\": false, \"matching_id\": null, \"mark_price\": 21758.49, \"liquidity\": \"T\", \"instrument_name\": \"BTC-PERPETUAL\", \"index_price\": 21755.39, \"fee_currency\": \"BTC\", \"fee\": 0.00000184, \"direction\": \"buy\", \"api\": false, \"amount\": 80 } ]"),
			},
			expect: &models.UserTradesNotification{
				models.UserTrade{
					UnderlyingPrice: 0,
					TradeSeq:        0x4dfbf83,
					TradeID:         "119813642",
					Timestamp:       0x1833066fbd8,
					TickDirection:   0,
					State:           "filled",
					SelfTrade:       false,
					ReduceOnly:      false,
					ProfitLost:      0,
					Price:           21760.5,
					PostOnly:        false,
					OrderType:       "market",
					OrderID:         "14228823973",
					MatchingID:      nil,
					MarkPrice:       21758.49,
					Liquidity:       "T",
					Liquidation:     "",
					Label:           "",
					IV:              0,
					InstrumentName:  "BTC-PERPETUAL",
					IndexPrice:      21755.39,
					FeeCurrency:     "BTC",
					Fee:             4.6e-07,
					Direction:       "buy",
					Amount:          20,
					BlockTradeID:    "",
				},
				models.UserTrade{
					UnderlyingPrice: 0,
					TradeSeq:        0x4dfbf84,
					TradeID:         "119813643",
					Timestamp:       0x1833066fbd8,
					TickDirection:   1,
					State:           "filled",
					SelfTrade:       false,
					ReduceOnly:      false,
					ProfitLost:      0,
					Price:           21760.5,
					PostOnly:        false,
					OrderType:       "market",
					OrderID:         "14228823973",
					MatchingID:      nil,
					MarkPrice:       21758.49,
					Liquidity:       "T",
					Liquidation:     "",
					Label:           "",
					IV:              0,
					InstrumentName:  "BTC-PERPETUAL",
					IndexPrice:      21755.39,
					FeeCurrency:     "BTC",
					Fee:             1.84e-06,
					Direction:       "buy",
					Amount:          80,
					BlockTradeID:    "",
				},
			},
		},
		{
			req: &jsonrpc2.Request{
				Method: "subscription",
			},
			params: Event{
				Channel: "instrument.state.BTC",
				Data:    json.RawMessage("{\"timestamp\":1662970320027,\"state\":\"terminated\",\"instrument_name\":\"BTC-11SEP22-16000-P\"}"),
			},
			expect: &models.InstrumentChangeNotification{
				InstrumentName: "BTC-11SEP22-16000-P",
				State:          "terminated",
				Timestamp:      1662970320027,
			},
		},
	}

	client := newClient()
	err := client.Start()
	require.NoError(t, err)

	eventCh := make(chan interface{}, 100)
	listener := func(event interface{}) {
		eventCh <- event
	}

	for _, test := range tests {
		err = test.req.SetParams(test.params)
		require.NoError(t, err)

		client.On(test.params.Channel, listener)
		client.Handle(context.Background(), nil, test.req)
		if assert.Len(t, eventCh, 1) {
			event := <-eventCh
			assert.Equal(t, test.expect, event)
		}
		client.Off(test.params.Channel, listener)
	}
}

func float64Pointer(v float64) *float64 {
	return &v
}
