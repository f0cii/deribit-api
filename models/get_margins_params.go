package models

type GetMarginsParams struct {
	InstrumentName string  `json:"instrument_name"`
	Amount         float64 `json:"amount"`
	Price          float64 `json:"price"`
}
