package models

import "github.com/shopspring/decimal"

type TickerResponse struct {
	BestAskAmount   decimal.Decimal `json:"best_ask_amount"`
	BestAskPrice    decimal.Decimal `json:"best_ask_price"`
	BestBidAmount   decimal.Decimal `json:"best_bid_amount"`
	BestBidPrice    decimal.Decimal `json:"best_bid_price"`
	CurrentFunding  decimal.Decimal `json:"current_funding"`
	Funding8H       decimal.Decimal `json:"funding_8h"`
	IndexPrice      decimal.Decimal `json:"index_price"`
	InstrumentName  string          `json:"instrument_name"`
	LastPrice       decimal.Decimal `json:"last_price"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	MaxPrice        decimal.Decimal `json:"max_price"`
	MinPrice        decimal.Decimal `json:"min_price"`
	OpenInterest    decimal.Decimal `json:"open_interest"`
	SettlementPrice decimal.Decimal `json:"settlement_price"`
	State           string          `json:"state"`
	Stats           Stats           `json:"stats"`
	Timestamp       uint64          `json:"timestamp"`
}
