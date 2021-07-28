package models

type GetTradingviewChartDataParams struct {
	InstrumentName string `json:"instrument_name"`
	StartTimestamp uint64 `json:"start_timestamp"`
	EndTimestamp   uint64 `json:"end_timestamp"`
	Resolution     string `json:"resolution"`
}
