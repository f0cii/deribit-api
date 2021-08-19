package models

import "github.com/shopspring/decimal"

type BuyParams struct {
	InstrumentName string           `json:"instrument_name"`
	Amount         decimal.Decimal  `json:"amount"`
	Type           string           `json:"type,omitempty"`
	Label          string           `json:"label,omitempty"`
	Price          *decimal.Decimal `json:"price,omitempty"`
	TimeInForce    string           `json:"time_in_force,omitempty"`
	MaxShow        *decimal.Decimal `json:"max_show,omitempty"`
	PostOnly       bool             `json:"post_only,omitempty"`
	RejectPostOnly bool             `json:"reject_post_only,omitempty"`
	ReduceOnly     bool             `json:"reduce_only,omitempty"`
	TriggerPrice   *decimal.Decimal `json:"trigger_price,omitempty"`
	Trigger        string           `json:"trigger,omitempty"`
	Advanced       string           `json:"advanced,omitempty"`
	MMP            *bool            `json:"mmp,omitempty"`
}
