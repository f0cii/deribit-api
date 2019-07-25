package models

type Deposit struct {
	Address           string  `json:"address"`
	Amount            float64 `json:"amount"`
	Currency          string  `json:"currency"`
	ReceivedTimestamp int64   `json:"received_timestamp"`
	State             string  `json:"state"`
	TransactionID     string  `json:"transaction_id"`
	UpdatedTimestamp  int64   `json:"updated_timestamp"`
}
