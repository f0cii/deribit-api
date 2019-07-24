package models

type GetFundingChartDataResponse struct {
	CurrentInterest float64     `json:"current_interest"`
	Data            [][]float64 `json:"data"`
	IndexPrice      float64     `json:"index_price"`
	Interest8H      float64     `json:"interest_8h"`
	Max             float64     `json:"max"`
}
