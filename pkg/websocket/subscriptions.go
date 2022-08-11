package websocket

import (
	"strings"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

func (c *Client) emitEvent(logger *zap.SugaredLogger, event *Event, i interface{}) {
	err := jsoniter.Unmarshal(event.Data, i)
	if err != nil {
		logger.Errorw("Fail to unmarshal data", "error", err)
		return
	}

	c.Emit(event.Channel, i)
}

func (c *Client) subscriptionsProcess(event *Event) {
	logger := c.l.With("func", "subscriptionsProcess", "channel", event.Channel)

	notification := getNotificationFromChannel(event.Channel)
	if notification == nil {
		logger.Infow("Not supported channel")
	}

	c.emitEvent(logger, event, notification)
}

// nolint:cyclop
func getNotificationFromChannel(channel string) interface{} {
	parts := strings.Split(channel, ".")
	if len(parts) == 1 && parts[0] == "announcements" {
		return &models.AnnouncementsNotification{}
	}
	if len(parts) < 2 {
		return nil
	}

	switch parts[0] {
	case "book":
		return getOrderBookNotificationFromChannel(channel)
	case "deribit_price_index":
		return &models.DeribitPriceIndexNotification{}
	case "deribit_price_ranking":
		return &models.DeribitPriceRankingNotification{}
	case "estimated_expiration_price":
		return &models.EstimatedExpirationPriceNotification{}
	case "markprice":
		if parts[1] == "options" {
			return &models.MarkpriceOptionsNotification{}
		}
	case "perpetual":
		return &models.PerpetualNotification{}
	case "quote":
		return &models.QuoteNotification{}
	case "ticker":
		return &models.TickerNotification{}
	case "trades":
		return &models.TradesNotification{}
	case "user":
		return getUserNotificationFromChannelParts(parts)
	case "instrument":
		if parts[1] == "state" {
			return &models.InstrumentChangeNotification{}
		}
	}

	return nil
}

func getOrderBookNotificationFromChannel(channel string) interface{} {
	count := strings.Count(channel, ".")
	if count == 2 {
		// book.BTC-PERPETUAL.raw, book.BTC-PERPETUAL.100ms
		if strings.HasSuffix(channel, ".raw") {
			return &models.OrderBookRawNotification{}
		}
		return &models.OrderBookNotification{}
	}

	if count == 4 {
		// book.ETH-PERPETUAL.100.1.100ms, ...
		return &models.OrderBookGroupNotification{}
	}

	return nil
}

func getUserNotificationFromChannelParts(parts []string) interface{} {
	switch parts[1] {
	case "changes":
		return &models.UserChangesNotification{}
	case "orders":
		if parts[len(parts)-1] == "raw" {
			return &models.Order{}
		}
		return &models.UserOrderNotification{}
	case "portfolio":
		return &models.PortfolioNotification{}
	case "trades":
		return &models.UserTradesNotification{}
	default:
		return nil
	}
}
