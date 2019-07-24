package models

type GetUserTradesByCurrencyAndTimeParams struct {
	Currency       string `json:"currency"`
	Kind           string `json:"kind,omitempty"`
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	Count          int    `json:"count,omitempty"`
	IncludeOld     bool   `json:"include_old,omitempty"`
	Sorting        string `json:"sorting,omitempty"`
}
