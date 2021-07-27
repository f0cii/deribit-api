package models

type PerpetualNotification struct {
	Timestamp  uint64  `json:"timestamp"`
	Interest   float64 `json:"interest"`
	IndexPrice float64 `json:"index_price"`
}
