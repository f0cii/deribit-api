package models

type CancelWithdrawalParams struct {
	Currency string `json:"currency"`
	Id       int    `json:"id"`
}
