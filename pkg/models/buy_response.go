package models

type BuyResponse struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}
