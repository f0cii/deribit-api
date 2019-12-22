package deribit

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/frankrap/deribit-api/models"
	"log"
	"strings"
)

func (c *Client) subscriptionsProcess(event *Event) {
	if event.Channel == "announcements" {
		var notification models.AnnouncementsNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "book") {
		var notification models.OrderBookNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "deribit_price_index") {
		var notification models.DeribitPriceIndexNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "deribit_price_ranking") {
		var notification models.DeribitPriceRankingNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "estimated_expiration_price") {
		var notification models.EstimatedExpirationPriceNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "markprice.options") {
		var notification models.MarkpriceOptionsNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "perpetual") {
		var notification models.PerpetualNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "quote") {
		var notification models.QuoteNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "ticker") {
		var notification models.TickerNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "trades") {
		var notification models.TradesNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.changes") {
		var notification models.UserChangesNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.orders") {
		if string(event.Data)[0] == '{' {
			var notification models.UserOrderNotification
			var order models.Order
			err := jsoniter.Unmarshal(event.Data, &order)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			notification = append(notification, order)
			c.Emit(event.Channel, &notification)
		} else {
			var notification models.UserOrderNotification
			err := jsoniter.Unmarshal(event.Data, &notification)
			if err != nil {
				log.Printf("%v", err)
				return
			}
			c.Emit(event.Channel, &notification)
		}
	} else if strings.HasPrefix(event.Channel, "user.portfolio") {
		var notification models.PortfolioNotification
		err := jsoniter.Unmarshal(event.Data, &notification)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		c.Emit(event.Channel, &notification)
	} else if strings.HasPrefix(event.Channel, "user.trades") {
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
