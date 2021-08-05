package models

import "github.com/shopspring/decimal"

type PerpetualNotification struct {
	Timestamp  uint64          `json:"timestamp"`
	Interest   decimal.Decimal `json:"interest"`
	IndexPrice decimal.Decimal `json:"index_price"`
}
