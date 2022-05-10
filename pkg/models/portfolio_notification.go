package models

import "github.com/shopspring/decimal"

type PortfolioNotification struct {
	TotalPl                    decimal.Decimal `json:"total_pl"`
	SessionUpl                 decimal.Decimal `json:"session_upl"`
	SessionRpl                 decimal.Decimal `json:"session_rpl"`
	ProjectedMaintenanceMargin decimal.Decimal `json:"projected_maintenance_margin"`
	ProjectedInitialMargin     decimal.Decimal `json:"projected_initial_margin"`
	ProjectedDeltaTotal        decimal.Decimal `json:"projected_delta_total"`
	PortfolioMarginingEnabled  bool            `json:"portfolio_margining_enabled"`
	OptionsVega                decimal.Decimal `json:"options_vega"`
	OptionsValue               decimal.Decimal `json:"options_value"`
	OptionsTheta               decimal.Decimal `json:"options_theta"`
	OptionsSessionUpl          decimal.Decimal `json:"options_session_upl"`
	OptionsSessionRpl          decimal.Decimal `json:"options_session_rpl"`
	OptionsPl                  decimal.Decimal `json:"options_pl"`
	OptionsGamma               decimal.Decimal `json:"options_gamma"`
	OptionsDelta               decimal.Decimal `json:"options_delta"`
	MarginBalance              decimal.Decimal `json:"margin_balance"`
	MaintenanceMargin          decimal.Decimal `json:"maintenance_margin"`
	InitialMargin              decimal.Decimal `json:"initial_margin"`
	FuturesSessionUpl          decimal.Decimal `json:"futures_session_upl"`
	FuturesSessionRpl          decimal.Decimal `json:"futures_session_rpl"`
	FuturesPl                  decimal.Decimal `json:"futures_pl"`
	EstimatedLiquidationRatio  decimal.Decimal `json:"estimated_liquidation_ratio"`
	Equity                     decimal.Decimal `json:"equity"`
	DeltaTotal                 decimal.Decimal `json:"delta_total"`
	Currency                   string          `json:"currency"`
	Balance                    decimal.Decimal `json:"balance"`
	AvailableWithdrawalFunds   decimal.Decimal `json:"available_withdrawal_funds"`
	AvailableFunds             decimal.Decimal `json:"available_funds"`
}
