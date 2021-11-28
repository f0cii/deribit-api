package models

type GetSubaccountsDetailsParams struct {
	Currency       string `json:"currency"`
	WithOpenOrders bool   `json:"with_open_orders,omitempty"`
}
