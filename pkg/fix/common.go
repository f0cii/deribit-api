package fix

import (
	"bytes"
	"crypto/rand"
	"errors"
	"strconv"
	"strings"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/tag"
)

const (
	fixVersion = "FIX.4.4"
)

var ErrClosed = errors.New("connection is closed")

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func newOrderBookNotificationChannel(instrument string) string {
	return "book." + instrument
}

func newTradeNotificationChannel(instrument string) string {
	return "trades." + instrument
}

// nolint:funlen,cyclop
func decodeExecutionReport(msg *quickfix.Message) (order models.Order, err error) {
	status, err := getOrderStatus(msg)
	if err != nil {
		return order, err
	}

	if status == enum.OrdStatus_REJECTED {
		reason, err2 := getText(msg)
		if err2 == nil {
			err = errors.New(reason)
		} else {
			err = err2
		}
		return order, err
	}

	instrument, err := getSymbol(msg)
	if err != nil {
		return order, err
	}

	orderID, err := getOrderID(msg)
	if err != nil {
		return order, err
	}

	orderType, err := getOrdType(msg)
	if err != nil {
		return order, err
	}

	side, err := getSide(msg)
	if err != nil {
		return order, err
	}

	amount, err := getOrderQty(msg)
	if err != nil {
		return order, err
	}

	filledAmount, err := getCumQty(msg)
	if err != nil {
		return order, err
	}

	price, err := getPrice(msg)
	if err != nil {
		return order, err
	}

	avgPrice, err := getAvgPx(msg)
	if err != nil {
		return order, err
	}

	commission, err := getCommission(msg)
	if err != nil {
		return order, err
	}

	maxShow, err := getMaxShow(msg)
	if err != nil {
		return order, err
	}

	transactTime, err := getTransactTime(msg)
	if err != nil {
		return order, err
	}

	label, err := getDeribitLabel(msg)
	if err != nil {
		return order, err
	}

	var execInst string
	execInst, err = getExecInst(msg)
	if err != nil {
		return order, err
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
	if strings.Contains(execInst, string(enum.ExecInst_PARTICIPANT_DONT_INITIATE)) {
		order.PostOnly = true
	}
	if strings.Contains(execInst, string(enum.ExecInst_DO_NOT_INCREASE)) {
		order.ReduceOnly = true
	}

	return order, nil
}

func decodeOrderStatus(status enum.OrdStatus) string {
	switch status {
	case enum.OrdStatus_NEW, enum.OrdStatus_PARTIALLY_FILLED:
		return "open"
	case enum.OrdStatus_FILLED:
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

func decodeTimeInForce(timeInForce enum.TimeInForce) string {
	switch timeInForce {
	case enum.TimeInForce_GOOD_TILL_CANCEL:
		return "good_til_cancelled"
	case enum.TimeInForce_DAY:
		return "good_til_day"
	case enum.TimeInForce_FILL_OR_KILL:
		return "fill_or_kill"
	case enum.TimeInForce_IMMEDIATE_OR_CANCEL:
		return "immediate_or_cancel"
	default:
		return ""
	}
}

func copyMessage(msg *quickfix.Message) (*quickfix.Message, error) {
	out := quickfix.NewMessage()
	err := quickfix.ParseMessage(out, bytes.NewBufferString(msg.String()))
	if err != nil {
		return nil, err
	}
	return out, nil
}

func floatToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// nolint:cyclop
func getReqIDTagFromMsgType(msgType enum.MsgType) (quickfix.Tag, error) {
	switch msgType {
	case enum.MsgType_SECURITY_LIST:
		return tag.SecurityReqID, nil
	case enum.MsgType_MARKET_DATA_REQUEST:
		return tag.MDReqID, nil
	case enum.MsgType_MARKET_DATA_REQUEST_REJECT:
		return tag.MDReqID, nil
	case enum.MsgType_MARKET_DATA_SNAPSHOT_FULL_REFRESH:
		return tag.MDReqID, nil
	case enum.MsgType_MARKET_DATA_INCREMENTAL_REFRESH:
		return tag.MDReqID, nil
	case enum.MsgType_EXECUTION_REPORT:
		return tag.OrigClOrdID, nil
	case enum.MsgType_ORDER_CANCEL_REJECT:
		return tag.ClOrdID, nil
	case enum.MsgType_ORDER_MASS_CANCEL_REPORT:
		return tag.OrderID, nil
	case enum.MsgType_POSITION_REPORT:
		return tag.PosReqID, nil
	case enum.MsgType_USER_RESPONSE:
		return tag.UserRequestID, nil
	case enum.MsgType_SECURITY_STATUS:
		return tag.SecurityStatusReqID, nil
	default:
		return 0, errors.New("request id tag not found")
	}
}
