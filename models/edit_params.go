package models

import "github.com/shopspring/decimal"

type EditParams struct {
	OrderID        string           `json:"order_id"`
	Amount         decimal.Decimal  `json:"amount"`
	Price          *decimal.Decimal `json:"price,omitempty"`
	PostOnly       *bool            `json:"post_only,omitempty"`
	ReduceOnly     *bool            `json:"reduce_only,omitempty"`
	RejectPostOnly *bool            `json:"reject_post_only,omitempty"`
	Advanced       string           `json:"advanced,omitempty"`
	TriggerPrice   *decimal.Decimal `json:"trigger_price,omitempty"`
	MMP            *bool            `json:"mmp,omitempty"`
}
