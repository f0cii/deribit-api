# deribit-api
Go library for using the Deribit's v2 Websocket API.

V2 API Documentation: https://docs.deribit.com/v2/

### Example

```
package main

import (
	deribit "github.com/KyberNetwork/deribit-api"
	"github.com/KyberNetwork/deribit-api/models"
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
	client, err := deribit.New(cfg)
	if err != nil {
	    log.Printf("%v", err)
	    return
	}

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
	client.On("announcements", func(e *models.AnnouncementsNotification) {
    
    })
    client.On("book.ETH-PERPETUAL.100.1.100ms", func(e *models.OrderBookGroupNotification) {

    })
    client.On("book.BTC-PERPETUAL.100ms", func(e *models.OrderBookNotification) {

    })
    client.On("book.BTC-PERPETUAL.raw", func(e *models.OrderBookRawNotification) {

    })
    client.On("deribit_price_index.btc_usd", func(e *models.DeribitPriceIndexNotification) {

    })
    client.On("deribit_price_ranking.btc_usd", func(e *models.DeribitPriceRankingNotification) {

    })
    client.On("estimated_expiration_price.btc_usd", func(e *models.EstimatedExpirationPriceNotification) {

    })
    client.On("markprice.options.btc_usd", func(e *models.MarkpriceOptionsNotification) {

    })
    client.On("perpetual.BTC-PERPETUAL.raw", func(e *models.PerpetualNotification) {

    })
    client.On("quote.BTC-PERPETUAL", func(e *models.QuoteNotification) {

    })
    client.On("ticker.BTC-PERPETUAL.raw", func(e *models.TickerNotification) {

    })
    client.On("trades.BTC-PERPETUAL.raw", func(e *models.TradesNotification) {

    })

    client.On("user.changes.BTC-PERPETUAL.raw", func(e *models.UserChangesNotification) {

    })
    client.On("user.changes.future.BTC.raw", func(e *models.UserChangesNotification) {

    })
    client.On("user.orders.BTC-PERPETUAL.raw", func(e *models.UserOrderNotification) {

    })
    client.On("user.orders.future.BTC.100ms", func(e *models.UserOrderNotification) {

    })
    client.On("user.portfolio.btc", func(e *models.PortfolioNotification) {

    })
    client.On("user.trades.BTC-PERPETUAL.raw", func(e *models.UserTradesNotification) {

    })
    client.On("user.trades.future.BTC.100ms", func(e *models.UserTradesNotification) {

    })
    
    client.Subscribe([]string{
    	"announcements",
    	"book.BTC-PERPETUAL.none.10.100ms",	// none/1,2,5,10,25,100,250
    	"book.BTC-PERPETUAL.100ms",	// type: snapshot/change
    	"book.BTC-PERPETUAL.raw",
    	"deribit_price_index.btc_usd",
    	"deribit_price_ranking.btc_usd",
    	"estimated_expiration_price.btc_usd",
    	"markprice.options.btc_usd",
    	"perpetual.BTC-PERPETUAL.raw",
    	"quote.BTC-PERPETUAL",
    	"ticker.BTC-PERPETUAL.raw",
    	"trades.BTC-PERPETUAL.raw",
    	"user.changes.BTC-PERPETUAL.raw",
    	"user.changes.future.BTC.raw",
    	"user.orders.BTC-PERPETUAL.raw",
    	"user.orders.future.BTC.100ms",
    	"user.portfolio.btc",
    	"user.trades.BTC-PERPETUAL.raw",
    	"user.trades.future.BTC.100ms",
    })

	forever := make(chan bool)
	<- forever
}

```