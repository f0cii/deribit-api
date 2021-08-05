package models

import "github.com/shopspring/decimal"

type WithdrawalPriority struct {
	Name  string          `json:"name"`
	Value decimal.Decimal `json:"value"`
}

type Currency struct {
	CoinType             string               `json:"coin_type"`
	Currency             string               `json:"currency"`
	CurrencyLong         string               `json:"currency_long"`
	FeePrecision         int                  `json:"fee_precision"`
	MinConfirmations     int                  `json:"min_confirmations"`
	MinWithdrawalFee     decimal.Decimal      `json:"min_withdrawal_fee"`
	WithdrawalFee        decimal.Decimal      `json:"withdrawal_fee"`
	WithdrawalPriorities []WithdrawalPriority `json:"withdrawal_priorities"`
}
