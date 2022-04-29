package models

import "github.com/shopspring/decimal"

type DeribitPriceRanking struct {
	Weight     decimal.Decimal `json:"weight"`
	Timestamp  uint64          `json:"timestamp"`
	Price      decimal.Decimal `json:"price"`
	Identifier string          `json:"identifier"`
	Enabled    bool            `json:"enabled"`
}

type DeribitPriceRankingNotification []DeribitPriceRanking
