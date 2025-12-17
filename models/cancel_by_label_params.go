package models

type CancelByLabelParams struct {
	Label    string `json:"label"`
	Currency string `json:"currency,omitempty"`
}
