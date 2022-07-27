package models

type Instrument struct {
	TickSize             float64 `json:"tick_size"`
	TakerCommission      float64 `json:"taker_commission"`
	SettlementPeriod     string  `json:"settlement_period"`
	QuoteCurrency        string  `json:"quote_currency"`
	MinTradeAmount       float64 `json:"min_trade_amount"`
	MakerCommission      float64 `json:"maker_commission"`
	Leverage             int     `json:"leverage"`
	Kind                 string  `json:"kind"`
	IsActive             bool    `json:"is_active"`
	InstrumentID         uint32  `json:"instrument_id"`
	InstrumentName       string  `json:"instrument_name"`
	ExpirationTimestamp  uint64  `json:"expiration_timestamp"`
	CreationTimestamp    uint64  `json:"creation_timestamp"`
	ContractSize         float64 `json:"contract_size"`
	BaseCurrency         string  `json:"base_currency"`
	BlockTradeCommission float64 `json:"block_trade_commission"`
	OptionType           string  `json:"option_type"`
	Strike               float64 `json:"strike"`
}
