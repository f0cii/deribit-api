package models

type Trade struct {
	TradeSeq       int     `json:"trade_seq"`
	TradeID        string  `json:"trade_id"`
	Timestamp      int64   `json:"timestamp"`
	TickDirection  int     `json:"tick_direction"`
	Price          float64 `json:"price"`
	Iv             float64 `json:"iv"`
	InstrumentName string  `json:"instrument_name"`
	IndexPrice     float64 `json:"index_price"`
	Direction      string  `json:"direction"`
	Amount         float64 `json:"amount"`
}
