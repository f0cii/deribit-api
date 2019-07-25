package models

type GetTransfersResponse struct {
	Count int        `json:"count"`
	Data  []Transfer `json:"data"`
}
