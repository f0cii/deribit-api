package models

import (
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"
)

type OrderBookGroupNotification struct {
	Timestamp      int64               `json:"timestamp"`
	InstrumentName string              `json:"instrument_name"`
	ChangeID       int64               `json:"change_id"`
	Bids           [][]decimal.Decimal `json:"bids"` // [price, amount]
	Asks           [][]decimal.Decimal `json:"asks"` // [price, amount]
}

// OrderBookNotificationItem ...
// ["change",6947.0,82640.0]
// ["new",6942.5,6940.0]
// ["delete",6914.0,0.0]
type OrderBookNotificationItem struct {
	Action string          `json:"action"`
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
}

func (item *OrderBookNotificationItem) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&item.Action, &item.Price, &item.Amount}
	wantLen := len(tmp)
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	if g, e := len(tmp), wantLen; g != e {
		return fmt.Errorf("wrong number of fields in Order: %d != %d", g, e)
	}
	return nil
}

type OrderBookNotification struct {
	Type           string                      `json:"type"`
	Timestamp      int64                       `json:"timestamp"`
	InstrumentName string                      `json:"instrument_name"`
	PrevChangeID   int64                       `json:"prev_change_id"`
	ChangeID       int64                       `json:"change_id"`
	Bids           []OrderBookNotificationItem `json:"bids"` // [action, price, amount]
	Asks           []OrderBookNotificationItem `json:"asks"` // [action, price, amount]
}

type OrderBookRawNotification struct {
	Timestamp      int64                       `json:"timestamp"`
	InstrumentName string                      `json:"instrument_name"`
	PrevChangeID   int64                       `json:"prev_change_id"`
	ChangeID       int64                       `json:"change_id"`
	Bids           []OrderBookNotificationItem `json:"bids"` // [action, price, amount]
	Asks           []OrderBookNotificationItem `json:"asks"` // [action, price, amount]
}
