package models

import "github.com/shopspring/decimal"

type GetOrderBookResponse struct {
	Timestamp       uint64              `json:"timestamp"`
	Stats           Stats               `json:"stats"`
	State           string              `json:"state"`
	SettlementPrice decimal.Decimal     `json:"settlement_price"`
	MinPrice        decimal.Decimal     `json:"min_price"`
	MaxPrice        decimal.Decimal     `json:"max_price"`
	MarkPrice       decimal.Decimal     `json:"mark_price"`
	MarkIV          decimal.Decimal     `json:"mark_iv"`
	LastPrice       decimal.Decimal     `json:"last_price"`
	InstrumentName  string              `json:"instrument_name"`
	IndexPrice      decimal.Decimal     `json:"index_price"`
	Funding8H       decimal.Decimal     `json:"funding_8h"`
	CurrentFunding  decimal.Decimal     `json:"current_funding"`
	ChangeID        uint64              `json:"change_id"`
	BidIV           decimal.Decimal     `json:"bid_iv"`
	Bids            [][]decimal.Decimal `json:"bids"`
	BestBidPrice    decimal.Decimal     `json:"best_bid_price"`
	BestBidAmount   decimal.Decimal     `json:"best_bid_amount"`
	BestAskPrice    decimal.Decimal     `json:"best_ask_price"`
	BestAskAmount   decimal.Decimal     `json:"best_ask_amount"`
	AskIV           decimal.Decimal     `json:"ask_iv"`
	Asks            [][]decimal.Decimal `json:"asks"`
	DeliveryPrice   decimal.Decimal     `json:"delivery_price"`
	Greeks          Greeks              `json:"greeks"`
	InterestRate    decimal.Decimal     `json:"interest_rate"`
	OpenInterest    decimal.Decimal     `json:"open_interest"`
	UnderlyingIndex decimal.Decimal     `json:"underlying_index"`
	UnderlyingPrice decimal.Decimal     `json:"underlying_price"`
}

type Greeks struct {
	Delta decimal.Decimal `json:"delta"`
	Gamma decimal.Decimal `json:"gamma"`
	RHO   decimal.Decimal `json:"rho"`
	Theta decimal.Decimal `json:"theta"`
	Vega  decimal.Decimal `json:"vega"`
}
