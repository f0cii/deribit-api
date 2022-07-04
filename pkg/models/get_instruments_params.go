package models

type GetInstrumentsParams struct {
	Currency string `json:"currency"`
	Kind     string `json:"kind,omitempty"`
	Expired  *bool  `json:"expired,omitempty"`
}

type GetInstrumentParams struct {
	InstrumentName string `json:"instrument_name"`
}
