package models

type UserChangesNotification struct {
	Trades    []UserTrade `json:"trades"`
	Positions []Position  `json:"positions"`
	Orders    []Order     `json:"orders"`
}
