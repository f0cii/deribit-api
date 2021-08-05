package models

import "github.com/shopspring/decimal"

type TradeVolume struct {
	PutsVolume       decimal.Decimal `json:"puts_volume"`
	PutsVolume30D    decimal.Decimal `json:"puts_volume_30d"`
	PutsVolume7D     decimal.Decimal `json:"puts_volume_7d"`
	FuturesVolume    decimal.Decimal `json:"futures_volume"`
	FuturesVolume30D decimal.Decimal `json:"futures_volume_30d"`
	FuturesVolume7D  decimal.Decimal `json:"futures_volume_7d"`
	CurrencyPair     string          `json:"currency_pair"`
	CallsVolume      decimal.Decimal `json:"calls_volume"`
	CallsVolume30D   decimal.Decimal `json:"calls_volume_30d"`
	CallsVolume7D    decimal.Decimal `json:"calls_volume_7d"`
}

type GetTradeVolumesResponse []TradeVolume
