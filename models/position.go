package models

import "github.com/shopspring/decimal"

type Position struct {
	AveragePrice              decimal.Decimal `json:"average_price"`
	AveragePriceUSD           decimal.Decimal `json:"average_price_usd"`
	Delta                     decimal.Decimal `json:"delta"`
	Direction                 string          `json:"direction"`
	EstimatedLiquidationPrice decimal.Decimal `json:"estimated_liquidation_price"`
	FloatingProfitLoss        decimal.Decimal `json:"floating_profit_loss"`
	FloatingProfitLossUSD     decimal.Decimal `json:"floating_profit_loss_usd"`
	Gamma                     decimal.Decimal `json:"gamma"`
	IndexPrice                decimal.Decimal `json:"index_price"`
	InitialMargin             decimal.Decimal `json:"initial_margin"`
	InstrumentName            string          `json:"instrument_name"`
	Kind                      string          `json:"kind"`
	Leverage                  int             `json:"leverage"`
	MaintenanceMargin         decimal.Decimal `json:"maintenance_margin"`
	MarkPrice                 decimal.Decimal `json:"mark_price"`
	OpenOrdersMargin          decimal.Decimal `json:"open_orders_margin"`
	RealizedFunding           decimal.Decimal `json:"realized_funding"`
	RealizedProfitLoss        decimal.Decimal `json:"realized_profit_loss"`
	SettlementPrice           decimal.Decimal `json:"settlement_price"`
	Size                      decimal.Decimal `json:"size"`
	SizeCurrency              decimal.Decimal `json:"size_currency"`
	Theta                     decimal.Decimal `json:"theta"`
	TotalProfitLoss           decimal.Decimal `json:"total_profit_loss"`
	Vega                      decimal.Decimal `json:"vega"`
}
