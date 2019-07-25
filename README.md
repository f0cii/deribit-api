# deribit-api
Go library for using the Deribit's v2 Websocket API.

V2 API Documentation: https://docs.deribit.com/v2/

### Example

```
package main

import (
	"github.com/sumorf/deribit-api"
	"github.com/sumorf/deribit-api/models"
	"log"
)

func main() {
	cfg := &deribit.Configuration{
		Addr:          deribit.TestBaseURL,
		ApiKey:        "AsJTU16U",
		SecretKey:     "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		AutoReconnect: true,
		DebugMode:     true,
	}
	client := deribit.New(cfg)

	client.GetTime()
	client.Test()

	var err error

	// GetBookSummaryByCurrency
	getBookSummaryByCurrencyParams := &models.GetBookSummaryByCurrencyParams{
		Currency: "BTC",
		Kind:     "future",
	}
	var getBookSummaryByCurrencyResult []models.BookSummary
	getBookSummaryByCurrencyResult, err = client.GetBookSummaryByCurrency(getBookSummaryByCurrencyParams)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", getBookSummaryByCurrencyResult)

	// GetOrderBook
	getOrderBookParams := &models.GetOrderBookParams{
		InstrumentName: "BTC-PERPETUAL",
		Depth:          5,
	}
	var getOrderBookResult models.GetOrderBookResponse
	getOrderBookResult, err = client.GetOrderBook(getOrderBookParams)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", getOrderBookResult)

	// GetPosition
	getPositionParams := &models.GetPositionParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	var getPositionResult models.Position
	getPositionResult, err = client.GetPosition(getPositionParams)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", getPositionResult)

	// Buy
	guyParams := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		Price:          6000.0,
		Type:           "limit",
	}
	var buyResult models.BuyResponse
	buyResult, err = client.Buy(guyParams)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	log.Printf("%v", buyResult)

	// Subscribe
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

	client.Subscribe([]string{
		"book.ETH-PERPETUAL.100.1.100ms",
		"book.BTC-PERPETUAL.100ms",
		"quote.BTC-PERPETUAL",
		"ticker.BTC-PERPETUAL.raw",
		"trades.BTC-PERPETUAL.raw",
		"user.orders.BTC-PERPETUAL.raw",
		"user.portfolio.btc",
		"user.trades.BTC-PERPETUAL.raw",
		"user.trades.future.BTC.100ms",
	})

	forever := make(chan bool)
	<- forever
}

```