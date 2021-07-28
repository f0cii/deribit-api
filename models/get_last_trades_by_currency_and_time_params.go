package models

type GetLastTradesByCurrencyAndTimeParams struct {
	Currency       string `json:"currency"`
	Kind           string `json:"kind,omitempty"`
	StartTimestamp uint64 `json:"start_timestamp"`
	EndTimestamp   uint64 `json:"end_timestamp"`
	Count          int    `json:"count,omitempty"`
	IncludeOld     bool   `json:"include_old,omitempty"`
	Sorting        string `json:"sorting,omitempty"`
}
