package models

import (
	"fmt"
	"strconv"
	"strings"
)

type OrderBookGroupNotification struct {
	Timestamp      int64       `json:"timestamp"`
	InstrumentName string      `json:"instrument_name"`
	ChangeID       int64       `json:"change_id"`
	Bids           [][]float64 `json:"bids"` // [price, amount]
	Asks           [][]float64 `json:"asks"` // [price, amount]
}

// OrderBookNotificationItem ...
// ["change",6947.0,82640.0]
// ["new",6942.5,6940.0]
// ["delete",6914.0,0.0]
type OrderBookNotificationItem struct {
	Action string  `json:"action"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

func (item *OrderBookNotificationItem) UnmarshalJSON(b []byte) error {
	// b: ["new",59786.0,10.0]
	// log.Printf("b=%v", string(b))
	s := strings.TrimLeft(string(b), "[")
	s = strings.TrimRight(s, "]")
	l := strings.Split(s, ",")

	if len(l) != 3 {
		return fmt.Errorf("fail to UnmarshalJSON [%v]", string(b))
	}

	item.Action = strings.ReplaceAll(l[0], `"`, "")
	item.Price, _ = strconv.ParseFloat(l[1], 64)
	item.Amount, _ = strconv.ParseFloat(l[2], 64)

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
