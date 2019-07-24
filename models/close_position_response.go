package models

type ClosePositionResponse struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}
