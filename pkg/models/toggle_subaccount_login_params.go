package models

type ToggleSubaccountLoginParams struct {
	Sid   int    `json:"sid"`
	State string `json:"state"`
}
