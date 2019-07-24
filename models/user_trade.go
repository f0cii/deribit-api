package models

type UserTrade struct {
	TradeSeq       int         `json:"trade_seq"`
	TradeID        string      `json:"trade_id"`
	Timestamp      int64       `json:"timestamp"`
	TickDirection  int         `json:"tick_direction"`
	State          string      `json:"state"`
	SelfTrade      bool        `json:"self_trade"`
	Price          float64     `json:"price"`
	OrderType      string      `json:"order_type"`
	OrderID        string      `json:"order_id"`
	MatchingID     interface{} `json:"matching_id"`
	Liquidity      string      `json:"liquidity"`
	InstrumentName string      `json:"instrument_name"`
	IndexPrice     float64     `json:"index_price"`
	FeeCurrency    string      `json:"fee_currency"`
	Fee            float64     `json:"fee"`
	Direction      string      `json:"direction"`
	Amount         float64     `json:"amount"`
}
