package models

import "github.com/shopspring/decimal"

type GetIndexResponse struct {
	BTC decimal.Decimal `json:"BTC"`
	ETH decimal.Decimal `json:"ETH"`
	Edp decimal.Decimal `json:"edp"`
}
