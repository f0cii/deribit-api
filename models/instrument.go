package models

import "github.com/shopspring/decimal"

type Instrument struct {
	TickSize             decimal.Decimal `json:"tick_size"`
	TakerCommission      decimal.Decimal `json:"taker_commission"`
	SettlementPeriod     string          `json:"settlement_period"`
	QuoteCurrency        string          `json:"quote_currency"`
	MinTradeAmount       decimal.Decimal `json:"min_trade_amount"`
	MakerCommission      decimal.Decimal `json:"maker_commission"`
	Leverage             int             `json:"leverage"`
	Kind                 string          `json:"kind"`
	IsActive             bool            `json:"is_active"`
	InstrumentName       string          `json:"instrument_name"`
	ExpirationTimestamp  uint64          `json:"expiration_timestamp"`
	CreationTimestamp    uint64          `json:"creation_timestamp"`
	ContractSize         decimal.Decimal `json:"contract_size"`
	BaseCurrency         string          `json:"base_currency"`
	BlockTradeCommission decimal.Decimal `json:"block_trade_commission"`
	OptionType           string          `json:"option_type"`
	Strike               decimal.Decimal `json:"strike"`
}
