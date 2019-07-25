package models

type WithdrawParams struct {
	Currency string  `json:"currency"`
	Address  string  `json:"address"`
	Amount   float64 `json:"amount"`
	Priority string  `json:"priority,omitempty"`
	Tfa      string  `json:"tfa,omitempty"`
}
