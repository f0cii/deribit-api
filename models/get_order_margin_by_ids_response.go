package models

import "github.com/shopspring/decimal"

type OrderMargin struct {
	InitialMargin decimal.Decimal `json:"initial_margin"`
	OrderID       string          `json:"order_id"`
}

type GetOrderMarginByIDsResponse []OrderMargin
