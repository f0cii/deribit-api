package models

type OrderMargin struct {
	InitialMargin float64 `json:"initial_margin"`
	OrderID       string  `json:"order_id"`
}

type GetOrderMarginByIDsResponse []OrderMargin
