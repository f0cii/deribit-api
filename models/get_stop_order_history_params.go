package models

type GetStopOrderHistoryParams struct {
	Currency       string `json:"currency"`
	InstrumentName string `json:"instrument_name"`
	Count          int    `json:"count,omitempty"`
	Continuation   string `json:"continuation,omitempty"`
}
