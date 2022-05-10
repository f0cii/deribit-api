package models

import "github.com/shopspring/decimal"

type FundingChartData struct {
	IndexPrice decimal.Decimal `json:"index_price"`
	Interest8H decimal.Decimal `json:"interest_8h"`
	Timestamp  uint64          `json:"timestamp"`
}
type GetFundingChartDataResponse struct {
	CurrentInterest decimal.Decimal    `json:"current_interest"`
	Data            []FundingChartData `json:"data"`
	Interest8H      decimal.Decimal    `json:"interest_8h"`
}
