package models

type TradeVolume struct {
	PutsVolume    float64 `json:"puts_volume"`
	FuturesVolume float64 `json:"futures_volume"`
	CurrencyPair  string  `json:"currency_pair"`
	CallsVolume   float64 `json:"calls_volume"`
}

type GetTradeVolumesResponse []TradeVolume
