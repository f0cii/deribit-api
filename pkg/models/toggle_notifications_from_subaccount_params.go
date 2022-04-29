package models

type ToggleNotificationsFromSubaccountParams struct {
	Sid   int  `json:"sid"`
	State bool `json:"state"`
}
