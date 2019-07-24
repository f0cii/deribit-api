package models

type GetSettlementHistoryByCurrencyParams struct {
	Currency string `json:"currency"`
	Type     string `json:"type,omitempty"`
	Count    int    `json:"count,omitempty"`
}
