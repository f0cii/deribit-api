package models

type GetUserTradesResponse struct {
	Trades  []UserTrade `json:"trades"`
	HasMore bool        `json:"has_more"`
}
