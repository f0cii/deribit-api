package models

type GetSettlementHistoryByInstrumentParams struct {
	InstrumentName string `json:"instrument_name"`
	Type           string `json:"type,omitempty"`
	Count          int    `json:"count,omitempty"`
}
