package models

import "github.com/shopspring/decimal"

type MarkpriceOption struct {
	MarkPrice      decimal.Decimal `json:"mark_price"`
	Iv             decimal.Decimal `json:"iv"`
	InstrumentName string          `json:"instrument_name"`
}

type MarkpriceOptionsNotification []MarkpriceOption
