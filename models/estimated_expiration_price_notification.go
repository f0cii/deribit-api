package models

import "github.com/shopspring/decimal"

type EstimatedExpirationPriceNotification struct {
	Seconds     uint64          `json:"seconds"`
	Price       decimal.Decimal `json:"price"`
	IsEstimated bool            `json:"is_estimated"`
}
