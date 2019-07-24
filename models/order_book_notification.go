package models

type OrderBookNotification struct {
	Type           string          `json:"type"`
	Timestamp      int64           `json:"timestamp"`
	InstrumentName string          `json:"instrument_name"`
	ChangeID       int64           `json:"change_id"`
	Bids           [][]interface{} `json:"bids"`
	Asks           [][]interface{} `json:"asks"`
}
