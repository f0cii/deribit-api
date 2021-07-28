package models

type Withdrawal struct {
	Address            string  `json:"address"`
	Amount             float64 `json:"amount"`
	ConfirmedTimestamp *uint64 `json:"confirmed_timestamp,omitempty"`
	CreatedTimestamp   uint64  `json:"created_timestamp"`
	Currency           string  `json:"currency"`
	Fee                float64 `json:"fee"`
	ID                 uint64  `json:"id"`
	Priority           float64 `json:"priority"`
	State              string  `json:"state"`
	TransactionID      *string `json:"transaction_id,omitempty"`
	UpdatedTimestamp   uint64  `json:"updated_timestamp"`
}
