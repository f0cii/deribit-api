package models

type EditParams struct {
	OrderID        string   `json:"order_id"`
	Amount         float64  `json:"amount"`
	Label          string   `json:"label,omitempty"`
	Price          *float64 `json:"price,omitempty"`
	PostOnly       *bool    `json:"post_only,omitempty"`
	ReduceOnly     *bool    `json:"reduce_only,omitempty"`
	RejectPostOnly *bool    `json:"reject_post_only,omitempty"`
	Advanced       string   `json:"advanced,omitempty"`
	TriggerPrice   *float64 `json:"trigger_price,omitempty"`
	MMP            *bool    `json:"mmp,omitempty"`
}
