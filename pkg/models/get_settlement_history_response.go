package models

type GetSettlementHistoryResponse struct {
	Settlements  []Settlement `json:"settlements"`
	Continuation string       `json:"continuation"`
}
