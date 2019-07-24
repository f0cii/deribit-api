package models

type EditParams struct {
	OrderID   string  `json:"order_id"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	PostOnly  bool    `json:"post_only,omitempty"`
	Advanced  string  `json:"advanced,omitempty"`
	StopPrice float64 `json:"stop_price,omitempty"`
}
