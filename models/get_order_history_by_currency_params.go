package models

type GetOrderHistoryByCurrencyParams struct {
	Currency        string `json:"currency"`
	Kind            string `json:"kind,omitempty"`
	Count           int    `json:"count,omitempty"`
	Offset          int    `json:"offset,omitempty"`
	IncludeOld      bool   `json:"include_old,omitempty"`
	IncludeUnfilled bool   `json:"include_unfilled,omitempty"`
}
