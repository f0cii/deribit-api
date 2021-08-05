package models

import "github.com/shopspring/decimal"

type Withdrawal struct {
	Address            string          `json:"address"`
	Amount             decimal.Decimal `json:"amount"`
	ConfirmedTimestamp *uint64         `json:"confirmed_timestamp,omitempty"`
	CreatedTimestamp   uint64          `json:"created_timestamp"`
	Currency           string          `json:"currency"`
	Fee                decimal.Decimal `json:"fee"`
	ID                 uint64          `json:"id"`
	Priority           decimal.Decimal `json:"priority"`
	State              string          `json:"state"`
	TransactionID      *string         `json:"transaction_id,omitempty"`
	UpdatedTimestamp   uint64          `json:"updated_timestamp"`
}
