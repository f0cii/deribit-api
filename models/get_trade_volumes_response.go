package models

type TradeVolume struct {
	PutsVolume       float64 `json:"puts_volume"`
	PutsVolume30D    float64 `json:"puts_volume_30d"`
	PutsVolume7D     float64 `json:"puts_volume_7d"`
	FuturesVolume    float64 `json:"futures_volume"`
	FuturesVolume30D float64 `json:"futures_volume_30d"`
	FuturesVolume7D  float64 `json:"futures_volume_7d"`
	CurrencyPair     string  `json:"currency_pair"`
	CallsVolume      float64 `json:"calls_volume"`
	CallsVolume30D   float64 `json:"calls_volume_30d"`
	CallsVolume7D    float64 `json:"calls_volume_7d"`
}

type GetTradeVolumesResponse []TradeVolume
