package models

type GetLastTradesResponse struct {
	Trades  []Trade `json:"trades"`
	HasMore bool    `json:"has_more"`
}
