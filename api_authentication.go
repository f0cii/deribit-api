package deribit

import "github.com/frankrap/deribit-api/models"

func (c *Client) Auth(apiKey string, secretKey string) (err error) {
	params := models.ClientCredentialsParams{
		GrantType:    "client_credentials",
		ClientID:     apiKey,
		ClientSecret: secretKey,
	}
	var result models.AuthResponse
	err = c.Call("public/auth", params, &result)
	if err != nil {
		return
	}
	c.auth.token = result.AccessToken
	c.auth.refresh = result.RefreshToken
	return
}

func (c *Client) Logout() (err error) {
	var result = struct {
	}{}
	err = c.Call("public/auth", nil, &result)
	return
}
