package models

import "github.com/shopspring/decimal"

type Stats struct {
	Volume      decimal.Decimal  `json:"volume"`
	PriceChange *decimal.Decimal `json:"price_change"`
	Low         decimal.Decimal  `json:"low"`
	High        decimal.Decimal  `json:"high"`
}
type TickerNotification struct {
	Timestamp       uint64           `json:"timestamp"`
	Stats           Stats            `json:"stats"`
	State           string           `json:"state"`
	SettlementPrice decimal.Decimal  `json:"settlement_price"`
	OpenInterest    decimal.Decimal  `json:"open_interest"`
	MinPrice        decimal.Decimal  `json:"min_price"`
	MaxPrice        decimal.Decimal  `json:"max_price"`
	MarkPrice       decimal.Decimal  `json:"mark_price"`
	LastPrice       decimal.Decimal  `json:"last_price"`
	InstrumentName  string           `json:"instrument_name"`
	IndexPrice      decimal.Decimal  `json:"index_price"`
	Funding8H       decimal.Decimal  `json:"funding_8h"`
	CurrentFunding  decimal.Decimal  `json:"current_funding"`
	BestBidPrice    *decimal.Decimal `json:"best_bid_price"`
	BestBidAmount   decimal.Decimal  `json:"best_bid_amount"`
	BestAskPrice    *decimal.Decimal `json:"best_ask_price"`
	BestAskAmount   decimal.Decimal  `json:"best_ask_amount"`
}
