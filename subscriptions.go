package deribit

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sumorf/deribit-api/models"
	"log"
	"strings"
)

func (c *Client) subscriptionsProcess(event *Event) {
	if strings.HasPrefix(event.Channel, "book.") {
		//pCount := strings.Count(event.Channel, ".")
		//if pCount == 4 {
		//	// book.{instrument_name}.{group}.{depth}.{interval}
		//} else if pCount == 2 {
		//	// book.{instrument_name}.{interval}
		//}
		var notification models.OrderBookNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "quote.") {
		var notification models.QuoteNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "ticker.") {
		var notification models.TickerNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "trades.") {
		var notification models.TradesNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.orders.") {
		log.Printf("%v", string(event.Data))
		var notification models.UserOrderNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.portfolio.") {
		var notification models.PortfolioNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.trades.") {
		var notification models.UserTradesNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else {
		log.Printf("%v", string(event.Data))
	}
}
