package models

import "github.com/shopspring/decimal"

type UserTrade struct {
	UnderlyingPrice decimal.Decimal `json:"underlying_price"`
	TradeSeq        uint64          `json:"trade_seq"`
	TradeID         string          `json:"trade_id"`
	Timestamp       uint64          `json:"timestamp"`
	TickDirection   int             `json:"tick_direction"`
	State           string          `json:"state"`
	SelfTrade       bool            `json:"self_trade"`
	ReduceOnly      bool            `json:"reduce_only"`
	ProfitLost      decimal.Decimal `json:"profit_lost"`
	Price           decimal.Decimal `json:"price"`
	PostOnly        bool            `json:"post_only"`
	OrderType       string          `json:"order_type"`
	OrderID         string          `json:"order_id"`
	MatchingID      *string         `json:"matching_id"`
	MarkPrice       decimal.Decimal `json:"mark_price"`
	Liquidity       string          `json:"liquidity"`
	Liquidation     string          `json:"liquidation"`
	Label           string          `json:"label"`
	IV              decimal.Decimal `json:"iv"`
	InstrumentName  string          `json:"instrument_name"`
	IndexPrice      decimal.Decimal `json:"index_price"`
	FeeCurrency     string          `json:"fee_currency"`
	Fee             decimal.Decimal `json:"fee"`
	Direction       string          `json:"direction"`
	Amount          decimal.Decimal `json:"amount"`
	BlockTradeID    string          `json:"block_trade_id"`
}
