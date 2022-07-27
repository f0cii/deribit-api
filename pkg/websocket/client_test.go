package websocket

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newClient() *Client {
	cfg := &Configuration{
		Addr:          TestBaseURL,
		ApiKey:        "AsJTU16U",
		SecretKey:     "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		AutoReconnect: true,
		DebugMode:     true,
	}
	c := New(zap.NewExample().Sugar(), cfg)
	_ = c.Start()
	return c
}

func TestClient_GetTime(t *testing.T) {
	client := newClient()
	tm, err := client.GetTime(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", tm)
}

func TestClient_Test(t *testing.T) {
	client := newClient()
	result, err := client.Test(context.Background())
	assert.Nil(t, err)
	t.Logf("%v", result)
}

func TestClient_GetBookSummaryByCurrency(t *testing.T) {
	client := newClient()
	params := &models.GetBookSummaryByCurrencyParams{
		Currency: "BTC",
		Kind:     "future",
	}
	result, err := client.GetBookSummaryByCurrency(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetBookSummaryByInstrument(t *testing.T) {
	client := newClient()
	params := &models.GetBookSummaryByInstrumentParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.GetBookSummaryByInstrument(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetOrderBook(t *testing.T) {
	client := newClient()
	params := &models.GetOrderBookParams{
		InstrumentName: "BTC-PERPETUAL",
		Depth:          5,
	}
	result, err := client.GetOrderBook(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_Ticker(t *testing.T) {
	client := newClient()
	params := &models.TickerParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.Ticker(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_GetPosition(t *testing.T) {
	client := newClient()
	params := &models.GetPositionParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	result, err := client.GetPosition(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_BuyMarket(t *testing.T) {
	client := newClient()
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         10,
		Type:           "market",
	}
	result, err := client.Buy(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_Buy(t *testing.T) {
	client := newClient()
	price := 6000.0
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40.0,
		Price:          &price,
		Type:           models.OrderTypeLimit,
	}
	result, err := client.Buy(context.Background(), params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestJsonOmitempty(t *testing.T) {
	maxShow := 40.0
	price := 6000.0
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40.0,
		Price:          &price,
		Type:           "limit",
		TimeInForce:    "good_til_cancelled",
		MaxShow:        &maxShow,
	}
	data, _ := json.Marshal(params)
	t.Log(string(data))
}
