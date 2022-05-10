package models

type GetOpenOrdersByCurrencyParams struct {
	Currency string `json:"currency"`
	Kind     string `json:"kind,omitempty"`
	Type     string `json:"type,omitempty"`
}
