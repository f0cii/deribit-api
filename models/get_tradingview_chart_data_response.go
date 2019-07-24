package models

type GetTradingviewChartDataResponse struct {
	Volume []float64 `json:"volume"`
	Ticks  []int64   `json:"ticks"`
	Status string    `json:"status"`
	Open   []float64 `json:"open"`
	Low    []float64 `json:"low"`
	High   []float64 `json:"high"`
	Close  []float64 `json:"close"`
}
