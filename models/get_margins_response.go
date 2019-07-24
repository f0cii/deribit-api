package models

type GetMarginsResponse struct {
	Buy      float64 `json:"buy"`
	MaxPrice float64 `json:"max_price"`
	MinPrice float64 `json:"min_price"`
	Sell     float64 `json:"sell"`
}
