package models

type Deposit struct {
	Address           string  `json:"address"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	ReceivedTimestamp uint64  `json:"received_timestamp"`
	State             string  `json:"state"`
	TransactionID     string  `json:"transaction_id"`
	UpdatedTimestamp  uint64  `json:"updated_timestamp"`
}
