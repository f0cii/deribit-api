package models

import "github.com/shopspring/decimal"

type QuoteNotification struct {
	Timestamp      uint64           `json:"timestamp"`
	InstrumentName string           `json:"instrument_name"`
	BestBidPrice   *decimal.Decimal `json:"best_bid_price"`
	BestBidAmount  decimal.Decimal  `json:"best_bid_amount"`
	BestAskPrice   *decimal.Decimal `json:"best_ask_price"`
	BestAskAmount  decimal.Decimal  `json:"best_ask_amount"`
}
