package models

type Settlement struct {
	Type              string  `json:"type"`
	Timestamp         int64   `json:"timestamp"`
	SessionProfitLoss float64 `json:"session_profit_loss"`
	ProfitLoss        float64 `json:"profit_loss"`
	Position          float64 `json:"position"`
	MarkPrice         float64 `json:"mark_price"`
	InstrumentName    string  `json:"instrument_name"`
	IndexPrice        float64 `json:"index_price"`
}
