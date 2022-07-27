package models

type Position struct {
	AveragePrice              float64 `json:"average_price"`
	AveragePriceUSD           float64 `json:"average_price_usd"`
	Delta                     float64 `json:"delta"`
	Direction                 string  `json:"direction"`
	EstimatedLiquidationPrice float64 `json:"estimated_liquidation_price"`
	FloatingProfitLoss        float64 `json:"floating_profit_loss"`
	FloatingProfitLossUSD     float64 `json:"floating_profit_loss_usd"`
	Gamma                     float64 `json:"gamma"`
	IndexPrice                float64 `json:"index_price"`
	InitialMargin             float64 `json:"initial_margin"`
	InstrumentName            string  `json:"instrument_name"`
	Kind                      string  `json:"kind"`
	Leverage                  int     `json:"leverage"`
	MaintenanceMargin         float64 `json:"maintenance_margin"`
	MarkPrice                 float64 `json:"mark_price"`
	OpenOrdersMargin          float64 `json:"open_orders_margin"`
	RealizedFunding           float64 `json:"realized_funding"`
	RealizedProfitLoss        float64 `json:"realized_profit_loss"`
	SettlementPrice           float64 `json:"settlement_price"`
	Size                      float64 `json:"size"`
	SizeCurrency              float64 `json:"size_currency"`
	Theta                     float64 `json:"theta"`
	TotalProfitLoss           float64 `json:"total_profit_loss"`
	Vega                      float64 `json:"vega"`
}
