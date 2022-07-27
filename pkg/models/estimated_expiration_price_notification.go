package models

type EstimatedExpirationPriceNotification struct {
	Seconds     uint64  `json:"seconds"`
	Price       float64 `json:"price"`
	IsEstimated bool    `json:"is_estimated"`
}
