package models

type GetStopOrderHistoryResponse struct {
	Entries      []StopOrder `json:"entries"`
	Continuation string      `json:"continuation"`
}
