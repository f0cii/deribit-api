package deribit

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/sumorf/deribit-api/models"
	"log"
	"testing"
)

func newClient() *Client {
	cfg := &Configuration{
		Addr:          TestBaseURL,
		ApiKey:        "AsJTU16U",
		SecretKey:     "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		AutoReconnect: true,
		DebugMode:     true,
	}
	client := New(cfg)
	return client
}

func TestClient_GetTime(t *testing.T) {
	client := newClient()
	tm, err := client.GetTime()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", tm)
}

func TestClient_Test(t *testing.T) {
	client := newClient()
	result, err := client.Test()
	assert.Nil(t, err)
	t.Logf("%v", result)
}

func TestClient_GetBookSummaryByCurrency(t *testing.T) {
	client := newClient()
	params := &models.GetBookSummaryByCurrencyParams{
		Currency: "BTC",
		Kind:     "future",
	}
	result, err := client.GetBookSummaryByCurrency(params)
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
	result, err := client.GetBookSummaryByInstrument(params)
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
	result, err := client.GetOrderBook(params)
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
	result, err := client.Ticker(params)
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
	result, err := client.GetPosition(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestClient_Buy(t *testing.T) {
	client := newClient()
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		Price:          6000.0,
		Type:           "limit",
	}
	result, err := client.Buy(params)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", result)
}

func TestJsonOmitempty(t *testing.T) {
	params := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		//Price:          6000.0,
		Type:        "limit",
		TimeInForce: "good_til_cancelled",
		MaxShow:     Float64Pointer(40.0),
	}
	data, _ := json.Marshal(params)
	t.Log(string(data))
}

func TestClient_Subscribe(t *testing.T) {
	client := newClient()

	client.On("book.ETH-PERPETUAL.100.1.100ms", func(e *models.OrderBookNotification) {
		log.Printf("e: %v", *e)
	})
	client.On("book.BTC-PERPETUAL.100ms", func(e *models.OrderBookNotification) {
		log.Printf("e: %v", *e)
	})
	client.On("quote.BTC-PERPETUAL", func(e *models.QuoteNotification) {
		log.Printf("e: %v", *e)
	})
	client.On("ticker.BTC-PERPETUAL.raw", func(e *models.TickerNotification) {
		log.Printf("e: %v", *e)
	})
	client.On("trades.BTC-PERPETUAL.raw", func(e *models.TradesNotification) {
		log.Printf("e: %#v", *e)
	})

	client.On("user.orders.BTC-PERPETUAL.raw", func(e *models.UserOrderNotification) {
		log.Printf("e: %#v", e)
	})
	client.On("user.portfolio.btc", func(e *models.PortfolioNotification) {
		log.Printf("e: %#v", e)
	})
	client.On("user.trades.BTC-PERPETUAL.raw", func(e *models.UserTradesNotification) {
		log.Printf("e: %#v", e)
	})
	client.On("user.trades.future.BTC.100ms", func(e *models.UserTradesNotification) {
		log.Printf("e: %#v", e)
	})

	// book.ETH-PERPETUAL.100.1.100ms
	// book.BTC-PERPETUAL.100ms
	// quote.BTC-PERPETUAL
	// ticker.BTC-PERPETUAL.raw
	// trades.BTC-PERPETUAL.raw

	// user.orders.BTC-PERPETUAL.raw
	// user.portfolio.btc
	// user.trades.BTC-PERPETUAL.raw
	// user.trades.future.BTC.100ms

	client.Subscribe([]string{
		"book.BTC-PERPETUAL.100ms",
		"user.orders.BTC-PERPETUAL.raw",
	})

	select {}
}
