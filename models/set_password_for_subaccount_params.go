package models

type SetPasswordForSubaccountParams struct {
	Sid      int    `json:"sid"`
	Password string `json:"password"`
}
