package models

type CancelWithdrawalParams struct {
	Currency string `json:"currency"`
	ID       int64  `json:"id"`
}
