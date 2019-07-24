package models

type ClosePositionParams struct {
	InstrumentName string  `json:"instrument_name"`
	Type           string  `json:"type"`
	Price          float64 `json:"price,omitempty"`
}
