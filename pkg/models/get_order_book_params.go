package models

type GetOrderBookParams struct {
	InstrumentName string `json:"instrument_name"`
	Depth          int    `json:"depth,omitempty"`
}
