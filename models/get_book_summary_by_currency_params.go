package models

type GetBookSummaryByCurrencyParams struct {
	Currency string `json:"currency"`
	Kind     string `json:"kind,omitempty"`
}
