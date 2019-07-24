package models

type GetUserTradesByCurrencyParams struct {
	Currency   string `json:"currency"`
	Kind       string `json:"kind,omitempty"`
	StartID    string `json:"start_id,omitempty"`
	EndID      string `json:"end_id,omitempty"`
	Count      int    `json:"count,omitempty"`
	IncludeOld bool   `json:"include_old,omitempty"`
	Sorting    string `json:"sorting,omitempty"`
}
