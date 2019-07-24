package models

type GetLastSettlementsByInstrumentParams struct {
	InstrumentName       string `json:"instrument_name"`
	Type                 string `json:"type,omitempty"`
	Count                int    `json:"count,omitempty"`
	Continuation         string `json:"continuation,omitempty"`
	SearchStartTimestamp int    `json:"search_start_timestamp,omitempty"`
}
