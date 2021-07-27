package models

type UserChangesNotification struct {
	InstrumentName string      `json:"instrument_name"`
	Trades         []UserTrade `json:"trades"`
	Positions      []Position  `json:"positions"`
	Orders         []Order     `json:"orders"`
}
