package websocket

import (
	"context"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) Auth(ctx context.Context) (result models.AuthResponse, err error) {
	params := models.ClientCredentialsParams{
		GrantType:    "client_credentials",
		ClientID:     c.apiKey,
		ClientSecret: c.secretKey,
	}
	err = c.Call(ctx, "public/auth", params, &result)
	if err != nil {
		return
	}
	return
}

func (c *Client) Logout(ctx context.Context) (err error) {
	var result struct{}
	err = c.Call(ctx, "private/logout", nil, &result)
	return
}
