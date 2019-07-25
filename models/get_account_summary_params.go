package models

type GetAccountSummaryParams struct {
	Currency string `json:"currency"`
	Extended bool   `json:"extended,omitempty"`
}
