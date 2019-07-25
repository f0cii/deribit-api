package deribit

import "github.com/sumorf/deribit-api/models"

func (c *Client) Auth() (err error) {
	params := models.ClientCredentialsParams{
		GrantType:    "client_credentials",
		ClientID:     c.apiKey,
		ClientSecret: c.secretKey,
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
