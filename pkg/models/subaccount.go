package models

type Subaccount struct {
	Email        string `json:"email"`
	ID           int    `json:"id"`
	IsPassword   bool   `json:"is_password"`
	LoginEnabled bool   `json:"login_enabled"`
	Portfolio    struct {
		Eth Portfolio `json:"eth"`
		Btc Portfolio `json:"btc"`
	} `json:"portfolio"`
	ReceiveNotifications bool   `json:"receive_notifications"`
	SystemName           string `json:"system_name"`
	TfaEnabled           bool   `json:"tfa_enabled"`
	Type                 string `json:"type"`
	Username             string `json:"username"`
}
