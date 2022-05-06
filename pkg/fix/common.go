package fix

import (
	"crypto/rand"
	"errors"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
)

const (
	fixVersion = "FIX.4.4"
)

var (
	ErrClosed = errors.New("connection is closed")
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func newOrderBookNotificationChannel(instrument string) string {
	return "book." + instrument
}

func decodeExecutionReport(msg *quickfix.Message) (order models.Order, err error) {
	status, err := getOrderStatus(msg)
	if err != nil {
		return
	}

	if status == enum.OrdStatus_REJECTED {
		reason, err2 := getText(msg)
		if err2 == nil {
			err = errors.New(reason)
		} else {
			err = err2
		}
		return
	}

	instrument, err := getSymbol(msg)
	if err != nil {
		return
	}

	orderID, err := getOrderID(msg)
	if err != nil {
		return
	}

	orderType, err := getOrdType(msg)
	if err != nil {
		return
	}

	side, err := getSide(msg)
	if err != nil {
		return
	}

	amount, err := getOrderQty(msg)
	if err != nil {
		return
	}

	filledAmount, err := getCumQty(msg)
	if err != nil {
		return
	}

	price, err := getPrice(msg)
	if err != nil {
		return
	}

	avgPrice, err := getAvgPx(msg)
	if err != nil {
		return
	}

	commission, err := getCommission(msg)
	if err != nil {
		return
	}

	maxShow, err := getMaxShow(msg)
	if err != nil {
		return
	}

	transactTime, err := getTransactTime(msg)
	if err != nil {
		return
	}

	label, err := getDeribitLabel(msg)
	if err != nil {
		return
	}

	order.OrderState = decodeOrderStatus(status)
	order.MaxShow = maxShow
	order.API = true
	order.Amount = amount
	order.Web = false
	order.InstrumentName = instrument
	order.Price = price
	order.LastUpdateTimestamp = uint64(transactTime.UnixMilli())
	order.FilledAmount = filledAmount
	order.AveragePrice = avgPrice
	order.OrderID = orderID
	order.Commission = commission
	order.Label = label
	order.CreationTimestamp = uint64(transactTime.UnixMilli())
	order.Direction = decodeOrderSide(side)
	order.OrderType = decodeOrderType(orderType)

	return
}

func decodeOrderStatus(status enum.OrdStatus) string {
	switch status {
	case enum.OrdStatus_NEW:
		return "open"
	case enum.OrdStatus_PARTIALLY_FILLED, enum.OrdStatus_FILLED:
		return "filled"
	case enum.OrdStatus_REJECTED:
		return "rejected"
	case enum.OrdStatus_CANCELED:
		return "cancelled"
	default:
		return ""
	}
}

func decodeOrderSide(side enum.Side) string {
	switch side {
	case enum.Side_BUY:
		return "buy"
	case enum.Side_SELL:
		return "sell"
	default:
		return ""
	}
}

func decodeOrderType(orderType enum.OrdType) string {
	switch orderType {
	case enum.OrdType_MARKET:
		return "market"
	case enum.OrdType_LIMIT:
		return "limit"
	case orderTypeStopMarket:
		return "stop_market"
	case enum.OrdType_STOP_LIMIT:
		return "stop_limit"
	default:
		return ""
	}
}
