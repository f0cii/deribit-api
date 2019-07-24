package models

type Instrument struct {
	TickSize            float64 `json:"tick_size"`
	Strike              float64 `json:"strike"`
	SettlementPeriod    string  `json:"settlement_period"`
	QuoteCurrency       string  `json:"quote_currency"`
	OptionType          string  `json:"option_type"`
	MinTradeAmount      float64 `json:"min_trade_amount"`
	Kind                string  `json:"kind"`
	IsActive            bool    `json:"is_active"`
	InstrumentName      string  `json:"instrument_name"`
	ExpirationTimestamp int64   `json:"expiration_timestamp"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	ContractSize        float64 `json:"contract_size"`
	BaseCurrency        string  `json:"base_currency"`
}
