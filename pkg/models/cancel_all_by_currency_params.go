package models

type CancelAllByCurrencyParams struct {
	Currency string `json:"currency"`
	Kind     string `json:"kind,omitempty"`
	Type     string `json:"type,omitempty"`
}
