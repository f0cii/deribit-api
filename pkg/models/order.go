package models

import "github.com/shopspring/decimal"

type Order struct {
	MMPCancelled        bool            `json:"mmp_cancelled"`
	OrderState          string          `json:"order_state"`
	MaxShow             decimal.Decimal `json:"max_show"`
	API                 bool            `json:"api"`
	Amount              decimal.Decimal `json:"amount"`
	Web                 bool            `json:"web"`
	InstrumentName      string          `json:"instrument_name"`
	Advanced            string          `json:"advanced,omitempty"`
	Triggered           *bool           `json:"triggered,omitempty"`
	BlockTrade          bool            `json:"block_trade"`
	OriginalOrderType   string          `json:"original_order_type"`
	Price               decimal.Decimal `json:"price"`
	TimeInForce         string          `json:"time_in_force"`
	AutoReplaced        bool            `json:"auto_replaced"`
	LastUpdateTimestamp uint64          `json:"last_update_timestamp"`
	PostOnly            bool            `json:"post_only"`
	Replaced            bool            `json:"replaced"`
	FilledAmount        decimal.Decimal `json:"filled_amount"`
	AveragePrice        decimal.Decimal `json:"average_price"`
	OrderID             string          `json:"order_id"`
	ReduceOnly          bool            `json:"reduce_only"`
	Commission          decimal.Decimal `json:"commission"`
	AppName             string          `json:"app_name"`
	Label               string          `json:"label"`
	TriggerOrderID      string          `json:"trigger_order_id"`
	TriggedPrice        decimal.Decimal `json:"trigger_price"`
	CreationTimestamp   uint64          `json:"creation_timestamp"`
	Direction           string          `json:"direction"`
	IsLiquidation       bool            `json:"is_liquidation"`
	OrderType           string          `json:"order_type"`
	USD                 decimal.Decimal `json:"usd,omitempty"`
	ProfitLoss          decimal.Decimal `json:"profit_loss"`
	Implv               decimal.Decimal `json:"implv,omitempty"`
	Trigger             string          `json:"trigger"`
}
