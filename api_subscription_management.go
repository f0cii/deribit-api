package deribit

import (
	"context"

	"github.com/KyberNetwork/deribit-api/models"
)

func (c *Client) PublicSubscribe(ctx context.Context, params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "public/subscribe", params, &result)
	return
}

func (c *Client) PublicUnsubscribe(ctx context.Context, params *models.UnsubscribeParams) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "public/unsubscribe", params, &result)
	return
}

func (c *Client) PrivateSubscribe(ctx context.Context, params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call(ctx, "private/subscribe", params, &result)
	return
}

func (c *Client) PrivateUnsubscribe(ctx context.Context, params *models.UnsubscribeParams) (result models.UnsubscribeResponse, err error) {
	err = c.Call(ctx, "private/unsubscribe", params, &result)
	return
}
