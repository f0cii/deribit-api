package models

import "github.com/shopspring/decimal"

type BookSummary struct {
	AskPrice               *decimal.Decimal `json:"ask_price"`
	BaseCurrency           string           `json:"base_currency"`
	BidPrice               *decimal.Decimal `json:"bid_price"`
	CreationTimestamp      uint64           `json:"creation_timestamp"`
	CurrentFunding         decimal.Decimal  `json:"current_funding"`
	EstimatedDeliveryPrice decimal.Decimal  `json:"estimated_delivery_price"`
	Funding8H              decimal.Decimal  `json:"funding_8h"`
	High                   decimal.Decimal  `json:"high"`
	InstrumentName         string           `json:"instrument_name"`
	InterestRate           decimal.Decimal  `json:"interest_rate"`
	Last                   *decimal.Decimal `json:"last"`
	Low                    *decimal.Decimal `json:"low"`
	MarkPrice              decimal.Decimal  `json:"mark_price"`
	MidPrice               *decimal.Decimal `json:"mid_price"`
	OpenInterest           decimal.Decimal  `json:"open_interest"`
	PriceChange            *decimal.Decimal `json:"price_change"`
	QuoteCurrency          string           `json:"quote_currency"`
	UnderlyingIndex        string           `json:"underlying_index"`
	UnderlyingPrice        decimal.Decimal  `json:"underlying_price"`
	VolumeUsd              decimal.Decimal  `json:"volume_usd"`
	Volume                 decimal.Decimal  `json:"volume"`
}
