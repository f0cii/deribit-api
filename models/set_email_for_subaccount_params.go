package models

type SetEmailForSubaccountParams struct {
	Sid   int    `json:"sid"`
	Email string `json:"email"`
}
