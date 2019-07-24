package models

type CancelAllByInstrumentParams struct {
	InstrumentName string `json:"instrument_name"`
	Type           string `json:"type,omitempty"`
}
