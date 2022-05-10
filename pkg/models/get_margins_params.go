package models

import "github.com/shopspring/decimal"

type GetMarginsParams struct {
	InstrumentName string          `json:"instrument_name"`
	Amount         decimal.Decimal `json:"amount"`
	Price          decimal.Decimal `json:"price"`
}
