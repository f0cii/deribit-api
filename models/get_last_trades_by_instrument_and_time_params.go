package models

type GetLastTradesByInstrumentAndTimeParams struct {
	InstrumentName string `json:"instrument_name"`
	StartTimestamp uint64 `json:"start_timestamp"`
	EndTimestamp   uint64 `json:"end_timestamp"`
	Count          int    `json:"count,omitempty"`
	IncludeOld     *bool  `json:"include_old,omitempty"`
	Sorting        string `json:"sorting,omitempty"`
}
