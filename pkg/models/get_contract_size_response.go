package models

import "github.com/shopspring/decimal"

type GetContractSizeResponse struct {
	ContractSize decimal.Decimal `json:"contract_size"`
}
