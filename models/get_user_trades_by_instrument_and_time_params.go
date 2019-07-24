package models

type GetUserTradesByInstrumentAndTimeParams struct {
	InstrumentName string `json:"instrument_name"`
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	Count          int    `json:"count,omitempty"`
	IncludeOld     bool   `json:"include_old,omitempty"`
	Sorting        string `json:"sorting,omitempty"`
}
