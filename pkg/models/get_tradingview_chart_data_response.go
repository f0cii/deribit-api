package models

import "github.com/shopspring/decimal"

type GetTradingviewChartDataResponse struct {
	Volume []decimal.Decimal `json:"volume"`
	Ticks  []uint64          `json:"ticks"`
	Status string            `json:"status"`
	Open   []decimal.Decimal `json:"open"`
	Low    []decimal.Decimal `json:"low"`
	High   []decimal.Decimal `json:"high"`
	Cost   []decimal.Decimal `json:"cost"`
	Close  []decimal.Decimal `json:"close"`
}
