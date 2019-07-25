package models

type CancelTransferByIdParams struct {
	Currency string `json:"currency"`
	Id       int    `json:"id"`
	Tfa      string `json:"tfa,omitempty"`
}
