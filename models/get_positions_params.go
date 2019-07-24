package models

type GetPositionsParams struct {
	Currency string `json:"currency"`
	Kind     string `json:"kind,omitempty"`
}
