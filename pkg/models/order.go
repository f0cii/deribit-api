package models

type Order struct {
	MMPCancelled        bool    `json:"mmp_cancelled"`
	OrderState          string  `json:"order_state"`
	MaxShow             float64 `json:"max_show"`
	API                 bool    `json:"api"`
	Amount              float64 `json:"amount"`
	Web                 bool    `json:"web"`
	InstrumentName      string  `json:"instrument_name"`
	Advanced            string  `json:"advanced,omitempty"`
	Triggered           *bool   `json:"triggered,omitempty"`
	BlockTrade          bool    `json:"block_trade"`
	OriginalOrderType   string  `json:"original_order_type"`
	Price               float64 `json:"price"`
	TimeInForce         string  `json:"time_in_force"`
	AutoReplaced        bool    `json:"auto_replaced"`
	LastUpdateTimestamp uint64  `json:"last_update_timestamp"`
	PostOnly            bool    `json:"post_only"`
	Replaced            bool    `json:"replaced"`
	FilledAmount        float64 `json:"filled_amount"`
	AveragePrice        float64 `json:"average_price"`
	OrderID             string  `json:"order_id"`
	ReduceOnly          bool    `json:"reduce_only"`
	Commission          float64 `json:"commission"`
	AppName             string  `json:"app_name"`
	Label               string  `json:"label"`
	TriggerOrderID      string  `json:"trigger_order_id"`
	TriggedPrice        float64 `json:"trigger_price"`
	CreationTimestamp   uint64  `json:"creation_timestamp"`
	Direction           string  `json:"direction"`
	IsLiquidation       bool    `json:"is_liquidation"`
	OrderType           string  `json:"order_type"`
	USD                 float64 `json:"usd,omitempty"`
	ProfitLoss          float64 `json:"profit_loss"`
	Implv               float64 `json:"implv,omitempty"`
	Trigger             string  `json:"trigger"`
}
