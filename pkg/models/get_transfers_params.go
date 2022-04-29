package models

type GetTransfersParams struct {
	Currency string `json:"currency"`
	Count    int    `json:"count,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}
