package models

import "github.com/shopspring/decimal"

type ClosePositionParams struct {
	InstrumentName string           `json:"instrument_name"`
	Type           string           `json:"type"`
	Price          *decimal.Decimal `json:"price,omitempty"`
}
