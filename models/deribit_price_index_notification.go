package models

type DeribitPriceIndexNotification struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
	IndexName string  `json:"index_name"`
}
