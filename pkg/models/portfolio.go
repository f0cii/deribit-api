package models

import "github.com/shopspring/decimal"

type Portfolio struct {
	AvailableFunds           decimal.Decimal `json:"available_funds"`
	AvailableWithdrawalFunds decimal.Decimal `json:"available_withdrawal_funds"`
	Balance                  decimal.Decimal `json:"balance"`
	Currency                 string          `json:"currency"`
	Equity                   decimal.Decimal `json:"equity"`
	InitialMargin            int             `json:"initial_margin"`
	MaintenanceMargin        int             `json:"maintenance_margin"`
	MarginBalance            decimal.Decimal `json:"margin_balance"`
}
