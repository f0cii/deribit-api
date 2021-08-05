package models

import "github.com/shopspring/decimal"

type Trade struct {
	Amount         decimal.Decimal `json:"amount"`
	BlockTradeID   string          `json:"block_trade_id"`
	Direction      string          `json:"direction"`
	IndexPrice     decimal.Decimal `json:"index_price"`
	InstrumentName string          `json:"instrument_name"`
	IV             decimal.Decimal `json:"iv"`
	Liquidation    string          `json:"liquidation"`
	MarkPrice      decimal.Decimal `json:"mark_price"`
	Price          decimal.Decimal `json:"price"`
	TickDirection  int             `json:"tick_direction"`
	Timestamp      uint64          `json:"timestamp"`
	TradeID        string          `json:"trade_id"`
	TradeSeq       uint64          `json:"trade_seq"`
}
