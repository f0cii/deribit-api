package models

type EditByLabelParams struct {
	Label          string   `json:"label"`
	InstrumentName string   `json:"instrument_name"`
	Amount         float64  `json:"amount"`
	Price          *float64 `json:"price,omitempty"`
	PostOnly       *bool    `json:"post_only,omitempty"`
	ReduceOnly     *bool    `json:"reduce_only,omitempty"`
	RejectPostOnly *bool    `json:"reject_post_only,omitempty"`
	Advanced       string   `json:"advanced,omitempty"`
	TriggerPrice   *float64 `json:"trigger_price,omitempty"`
	MMP            *bool    `json:"mmp,omitempty"`
}

type CancelAllByLabelParams struct {
	Label string `json:"label"`
}
