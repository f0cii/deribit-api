package models

import "github.com/shopspring/decimal"

type GetMarginsResponse struct {
	Buy      decimal.Decimal `json:"buy"`
	MaxPrice decimal.Decimal `json:"max_price"`
	MinPrice decimal.Decimal `json:"min_price"`
	Sell     decimal.Decimal `json:"sell"`
}
