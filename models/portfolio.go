package models

type Portfolio struct {
	AvailableFunds           int    `json:"available_funds"`
	AvailableWithdrawalFunds int    `json:"available_withdrawal_funds"`
	Balance                  int    `json:"balance"`
	Currency                 string `json:"currency"`
	Equity                   int    `json:"equity"`
	InitialMargin            int    `json:"initial_margin"`
	MaintenanceMargin        int    `json:"maintenance_margin"`
	MarginBalance            int    `json:"margin_balance"`
}
