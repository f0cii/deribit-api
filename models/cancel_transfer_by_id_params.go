package models

type CancelTransferByIDParams struct {
	Currency string `json:"currency"`
	ID       int    `json:"id"`
	Tfa      string `json:"tfa,omitempty"`
}
