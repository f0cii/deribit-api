package models

type GetOrderHistoryByInstrumentParams struct {
	InstrumentName  string `json:"instrument_name"`
	Count           int    `json:"count,omitempty"`
	Offset          int    `json:"offset,omitempty"`
	IncludeOld      bool   `json:"include_old,omitempty"`
	IncludeUnfilled bool   `json:"include_unfilled,omitempty"`
}
