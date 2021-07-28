package models

type GetLastSettlementsByCurrencyParams struct {
	Currency             string `json:"currency"`
	Type                 string `json:"type,omitempty"`
	Count                uint64 `json:"count,omitempty"`
	Continuation         string `json:"continuation,omitempty"`
	SearchStartTimestamp uint64 `json:"search_start_timestamp,omitempty"`
}
