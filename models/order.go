package models

import "strconv"

type Price float64

func (p *Price) UnmarshalJSON(data []byte) (err error) {
	if string(data) == `"market_price"` {
		*p = 0
		return
	}
	var f float64
	f, err = strconv.ParseFloat(string(data), 0)
	if err != nil {
		return
	}
	*p = Price(f)
	return
}

func (p *Price) ToFloat64() float64 {
	return float64(*p)
}

type Order struct {
	Advanced            string  `json:"advanced,omitempty"`
	Amount              float64 `json:"amount"`
	API                 bool    `json:"api"`
	TimeInForce         string  `json:"time_in_force"`
	ReduceOnly          bool    `json:"reduce_only"`
	ProfitLoss          float64 `json:"profit_loss"`
	Price               Price   `json:"price"`
	PostOnly            bool    `json:"post_only"`
	StopPrice           float64 `json:"stop_price,omitempty"`
	Triggered           bool    `json:"triggered,omitempty"`
	OrderType           string  `json:"order_type"`
	OrderState          string  `json:"order_state"`
	OrderID             string  `json:"order_id"`
	MaxShow             float64 `json:"max_show"`
	LastUpdateTimestamp int64   `json:"last_update_timestamp"`
	Label               string  `json:"label"`
	IsLiquidation       bool    `json:"is_liquidation"`
	InstrumentName      string  `json:"instrument_name"`
	FilledAmount        float64 `json:"filled_amount"`
	Direction           string  `json:"direction"`
	CreationTimestamp   int64   `json:"creation_timestamp"`
	Commission          float64 `json:"commission"`
	AveragePrice        float64 `json:"average_price"`
	Implv               float64 `json:"implv,omitempty"`
	Usd                 float64 `json:"usd,omitempty"`
}
