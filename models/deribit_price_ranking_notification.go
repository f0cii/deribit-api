package models

type DeribitPriceRanking struct {
	Weight     float64 `json:"weight"`
	Timestamp  int64   `json:"timestamp"`
	Price      float64 `json:"price"`
	Identifier string  `json:"identifier"`
	Enabled    bool    `json:"enabled"`
}

type DeribitPriceRankingNotification []DeribitPriceRanking
