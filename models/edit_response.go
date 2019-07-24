package models

type EditResponse struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}
