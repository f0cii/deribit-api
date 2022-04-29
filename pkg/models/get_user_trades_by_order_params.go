package models

type GetUserTradesByOrderParams struct {
	OrderID string `json:"order_id"`
	Sorting string `json:"sorting,omitempty"`
}
