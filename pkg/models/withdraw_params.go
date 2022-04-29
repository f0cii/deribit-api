package models

import "github.com/shopspring/decimal"

type WithdrawParams struct {
	Currency string          `json:"currency"`
	Address  string          `json:"address"`
	Amount   decimal.Decimal `json:"amount"`
	Priority string          `json:"priority,omitempty"`
	Tfa      string          `json:"tfa,omitempty"`
}
