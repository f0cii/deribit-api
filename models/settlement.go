package models

import "github.com/shopspring/decimal"

type Settlement struct {
	Funded            decimal.Decimal `json:"funded"`
	Funding           decimal.Decimal `json:"funding"`
	Type              string          `json:"type"`
	Timestamp         uint64          `json:"timestamp"`
	SessionProfitLoss decimal.Decimal `json:"session_profit_loss"`
	ProfitLoss        decimal.Decimal `json:"profit_loss"`
	Position          decimal.Decimal `json:"position"`
	MarkPrice         decimal.Decimal `json:"mark_price"`
	InstrumentName    string          `json:"instrument_name"`
	IndexPrice        decimal.Decimal `json:"index_price"`
	SessionBankrupcy  decimal.Decimal `json:"session_bankrupcy"`
	SessionTax        decimal.Decimal `json:"session_tax"`
	SessionTaxRate    decimal.Decimal `json:"session_tax_rate"`
	Socialized        decimal.Decimal `json:"socialized"`
}
