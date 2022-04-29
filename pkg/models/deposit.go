package models

import "github.com/shopspring/decimal"

type Deposit struct {
	Address           string          `json:"address"`
	Amount            decimal.Decimal `json:"amount"`
	Currency          string          `json:"currency"`
	ReceivedTimestamp uint64          `json:"received_timestamp"`
	State             string          `json:"state"`
	TransactionID     string          `json:"transaction_id"`
	UpdatedTimestamp  uint64          `json:"updated_timestamp"`
}
