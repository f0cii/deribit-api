package models

import "github.com/shopspring/decimal"

type Transfer struct {
	Amount           decimal.Decimal `json:"amount"`
	CreatedTimestamp uint64          `json:"created_timestamp"`
	Currency         string          `json:"currency"`
	Direction        string          `json:"direction"`
	ID               int             `json:"id"`
	OtherSide        string          `json:"other_side"`
	State            string          `json:"state"`
	Type             string          `json:"type"`
	UpdatedTimestamp uint64          `json:"updated_timestamp"`
}
