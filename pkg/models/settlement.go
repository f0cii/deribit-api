package models

type Settlement struct {
	Funded            float64 `json:"funded"`
	Funding           float64 `json:"funding"`
	Type              string  `json:"type"`
	Timestamp         uint64  `json:"timestamp"`
	SessionProfitLoss float64 `json:"session_profit_loss"`
	ProfitLoss        float64 `json:"profit_loss"`
	Position          float64 `json:"position"`
	MarkPrice         float64 `json:"mark_price"`
	InstrumentName    string  `json:"instrument_name"`
	IndexPrice        float64 `json:"index_price"`
	SessionBankrupcy  float64 `json:"session_bankrupcy"`
	SessionTax        float64 `json:"session_tax"`
	SessionTaxRate    float64 `json:"session_tax_rate"`
	Socialized        float64 `json:"socialized"`
}
