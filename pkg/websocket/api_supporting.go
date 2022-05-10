package websocket

import (
	"context"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) GetTime(ctx context.Context) (result int64, err error) {
	err = c.Call(ctx, "public/get_time", nil, &result)
	return
}

func (c *Client) Hello(ctx context.Context, params *models.HelloParams) (result models.HelloResponse, err error) {
	err = c.Call(ctx, "public/hello", params, &result)
	return
}

func (c *Client) Test(ctx context.Context) (result models.TestResponse, err error) {
	err = c.Call(ctx, "public/test", nil, &result)
	return
}
