package websocket

import (
	"context"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) SetHeartbeat(
	ctx context.Context,
	params *models.SetHeartbeatParams,
) (result string, err error) {
	err = c.Call(ctx, "public/set_heartbeat", params, &result)
	return
}

func (c *Client) DisableHeartbeat(ctx context.Context) (result string, err error) {
	err = c.Call(ctx, "public/disable_heartbeat", nil, &result)
	return
}

func (c *Client) EnableCancelOnDisconnect(
	ctx context.Context,
	params *models.SessionParams,
) (result string, err error) {
	err = c.Call(ctx, "private/enable_cancel_on_disconnect", params, &result)
	return
}

func (c *Client) DisableCancelOnDisconnect(
	ctx context.Context,
	params *models.SessionParams,
) (result string, err error) {
	err = c.Call(ctx, "private/disable_cancel_on_disconnect", params, &result)
	return
}
