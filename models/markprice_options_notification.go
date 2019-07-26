package models

type MarkpriceOption struct {
	MarkPrice      float64 `json:"mark_price"`
	Iv             float64 `json:"iv"`
	InstrumentName string  `json:"instrument_name"`
}

type MarkpriceOptionsNotification []MarkpriceOption
