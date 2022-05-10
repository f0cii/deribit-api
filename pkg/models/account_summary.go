package models

import "github.com/shopspring/decimal"

type AccountSummary struct {
	AvailableFunds            decimal.Decimal `json:"available_funds"`
	AvailableWithdrawalFunds  decimal.Decimal `json:"available_withdrawal_funds"`
	Balance                   decimal.Decimal `json:"balance"`
	Currency                  string          `json:"currency"`
	DeltaTotal                decimal.Decimal `json:"delta_total"`
	DepositAddress            string          `json:"deposit_address"`
	Email                     string          `json:"email"`
	Equity                    decimal.Decimal `json:"equity"`
	FuturesPl                 decimal.Decimal `json:"futures_pl"`
	FuturesSessionRpl         decimal.Decimal `json:"futures_session_rpl"`
	FuturesSessionUpl         decimal.Decimal `json:"futures_session_upl"`
	ID                        int             `json:"id"`
	InitialMargin             decimal.Decimal `json:"initial_margin"`
	MaintenanceMargin         decimal.Decimal `json:"maintenance_margin"`
	MarginBalance             decimal.Decimal `json:"margin_balance"`
	OptionsDelta              decimal.Decimal `json:"options_delta"`
	OptionsGamma              decimal.Decimal `json:"options_gamma"`
	OptionsPl                 decimal.Decimal `json:"options_pl"`
	OptionsSessionRpl         decimal.Decimal `json:"options_session_rpl"`
	OptionsSessionUpl         decimal.Decimal `json:"options_session_upl"`
	OptionsTheta              decimal.Decimal `json:"options_theta"`
	OptionsVega               decimal.Decimal `json:"options_vega"`
	PortfolioMarginingEnabled bool            `json:"portfolio_margining_enabled"`
	SessionFunding            decimal.Decimal `json:"session_funding"`
	SessionRpl                decimal.Decimal `json:"session_rpl"`
	SessionUpl                decimal.Decimal `json:"session_upl"`
	SystemName                string          `json:"system_name"`
	TfaEnabled                bool            `json:"tfa_enabled"`
	TotalPl                   decimal.Decimal `json:"total_pl"`
	Type                      string          `json:"type"`
	Username                  string          `json:"username"`
}
