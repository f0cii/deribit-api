package models

type PortfolioNotification struct {
	TotalPl                   float64 `json:"total_pl"`
	SessionUpl                float64 `json:"session_upl"`
	SessionRpl                float64 `json:"session_rpl"`
	SessionFunding            float64 `json:"session_funding"`
	PortfolioMarginingEnabled bool    `json:"portfolio_margining_enabled"`
	OptionsVega               float64 `json:"options_vega"`
	OptionsTheta              float64 `json:"options_theta"`
	OptionsSessionUpl         float64 `json:"options_session_upl"`
	OptionsSessionRpl         float64 `json:"options_session_rpl"`
	OptionsPl                 float64 `json:"options_pl"`
	OptionsGamma              float64 `json:"options_gamma"`
	OptionsDelta              float64 `json:"options_delta"`
	MarginBalance             float64 `json:"margin_balance"`
	MaintenanceMargin         float64 `json:"maintenance_margin"`
	InitialMargin             float64 `json:"initial_margin"`
	FuturesSessionUpl         float64 `json:"futures_session_upl"`
	FuturesSessionRpl         float64 `json:"futures_session_rpl"`
	FuturesPl                 float64 `json:"futures_pl"`
	Equity                    float64 `json:"equity"`
	DeltaTotal                float64 `json:"delta_total"`
	Currency                  string  `json:"currency"`
	Balance                   float64 `json:"balance"`
	AvailableWithdrawalFunds  float64 `json:"available_withdrawal_funds"`
	AvailableFunds            float64 `json:"available_funds"`
}
