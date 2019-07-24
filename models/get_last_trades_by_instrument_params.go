package models

type GetLastTradesByInstrumentParams struct {
	InstrumentName string `json:"instrument_name"`
	StartSeq       int    `json:"start_seq,omitempty"`
	EndSeq         int    `json:"end_seq,omitempty"`
	Count          int    `json:"count,omitempty"`
	IncludeOld     bool   `json:"include_old,omitempty"`
	Sorting        string `json:"sorting,omitempty"`
}
