package models

type FundingChartData struct {
	IndexPrice float64 `json:"index_price"`
	Interest8H float64 `json:"interest_8h"`
	Timestamp  uint64  `json:"timestamp"`
}
type GetFundingChartDataResponse struct {
	CurrentInterest float64            `json:"current_interest"`
	Data            []FundingChartData `json:"data"`
	Interest8H      float64            `json:"interest_8h"`
}
