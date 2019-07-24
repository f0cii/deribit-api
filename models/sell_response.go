package models

type SellResponse struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}
