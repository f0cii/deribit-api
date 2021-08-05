package models

type CancelTransferByIDParams struct {
	Currency string `json:"currency"`
	ID       int64  `json:"id"`
	Tfa      string `json:"tfa,omitempty"`
}
