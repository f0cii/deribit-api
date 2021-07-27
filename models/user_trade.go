package models

type UserTrade struct {
	UnderlyingPrice float64     `json:"underlying_price"`
	TradeSeq        uint64      `json:"trade_seq"`
	TradeID         string      `json:"trade_id"`
	Timestamp       uint64      `json:"timestamp"`
	TickDirection   int         `json:"tick_direction"`
	State           string      `json:"state"`
	SelfTrade       bool        `json:"self_trade"`
	ReduceOnly      bool        `json:"reduce_only"`
	ProfitLost      float64     `json:"profit_lost"`
	Price           float64     `json:"price"`
	PostOnly        bool        `json:"post_only"`
	OrderType       string      `json:"order_type"`
	OrderID         string      `json:"order_id"`
	MatchingID      interface{} `json:"matching_id"`
	MarkPrice       float64     `json:"mark_price"`
	Liquidity       string      `json:"liquidity"`
	Liquidation     string      `json:"liquidation"`
	Label           string      `json:"label"`
	IV              float64     `json:"iv"`
	InstrumentName  string      `json:"instrument_name"`
	IndexPrice      float64     `json:"index_price"`
	FeeCurrency     string      `json:"fee_currency"`
	Fee             float64     `json:"fee"`
	Direction       string      `json:"direction"`
	Amount          float64     `json:"amount"`
	BlockTradeID    string      `json:"block_trade_id"`
}
