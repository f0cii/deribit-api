package models

type GetDepositsResponse struct {
	Count int       `json:"count"`
	Data  []Deposit `json:"data"`
}
