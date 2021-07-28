package models

type BookSummary struct {
	AskPrice               *float64 `json:"ask_price"`
	BaseCurrency           string   `json:"base_currency"`
	BidPrice               *float64 `json:"bid_price"`
	CreationTimestamp      uint64   `json:"creation_timestamp"`
	CurrentFunding         float64  `json:"current_funding"`
	EstimatedDeliveryPrice float64  `json:"estimated_delivery_price"`
	Funding8H              float64  `json:"funding_8h"`
	High                   float64  `json:"high"`
	InstrumentName         string   `json:"instrument_name"`
	InterestRate           float64  `json:"interest_rate"`
	Last                   *float64 `json:"last"`
	Low                    *float64 `json:"low"`
	MarkPrice              float64  `json:"mark_price"`
	MidPrice               *float64 `json:"mid_price"`
	OpenInterest           float64  `json:"open_interest"`
	PriceChange            *float64 `json:"price_change"`
	QuoteCurrency          string   `json:"quote_currency"`
	UnderlyingIndex        string   `json:"underlying_index"`
	UnderlyingPrice        float64  `json:"underlying_price"`
	VolumeUsd              float64  `json:"volume_usd"`
	Volume                 float64  `json:"volume"`
}
