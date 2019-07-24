package models

type GetLastSettlementsResponse struct {
	Settlements  []Settlement `json:"settlements"`
	Continuation string       `json:"continuation"`
}
