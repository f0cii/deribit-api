package deribit

import "encoding/json"

var emptyParams = json.RawMessage("{}")

// privateParams is interface for methods require access_token
type privateParams interface {
	setToken(token string)
}

// Token is used to embedded in params for private methods
type Token struct {
	AccessToken string `json:"access_token"`
}

func (t *Token) setToken(token string) {
	t.AccessToken = token
}
