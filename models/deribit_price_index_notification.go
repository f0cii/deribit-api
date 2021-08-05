package models

import "github.com/shopspring/decimal"

type DeribitPriceIndexNotification struct {
	Timestamp int64           `json:"timestamp"`
	Price     decimal.Decimal `json:"price"`
	IndexName string          `json:"index_name"`
}
