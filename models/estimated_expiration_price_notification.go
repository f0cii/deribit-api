package models

type EstimatedExpirationPriceNotification struct {
	Seconds     int     `json:"seconds"`
	Price       float64 `json:"price"`
	IsEstimated bool    `json:"is_estimated"`
}
