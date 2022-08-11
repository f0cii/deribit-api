package websocket

import (
	"context"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) GetAnnouncements(
	ctx context.Context,
) (result []models.Announcement, err error) {
	err = c.Call(ctx, "public/get_announcements", nil, &result)
	return
}

func (c *Client) ChangeSubaccountName(
	ctx context.Context,
	params *models.ChangeSubaccountNameParams,
) (result string, err error) {
	err = c.Call(ctx, "private/change_subaccount_name", params, &result)
	return
}

func (c *Client) CreateSubaccount(ctx context.Context) (result models.Subaccount, err error) {
	err = c.Call(ctx, "private/create_subaccount", nil, &result)
	return
}

func (c *Client) DisableTfaForSubaccount(
	ctx context.Context,
	params *models.DisableTfaForSubaccountParams,
) (result string, err error) {
	err = c.Call(ctx, "private/disable_tfa_for_subaccount", params, &result)
	return
}

func (c *Client) GetAccountSummary(
	ctx context.Context,
	params *models.GetAccountSummaryParams,
) (result models.AccountSummary, err error) {
	err = c.Call(ctx, "private/get_account_summary", params, &result)
	return
}

func (c *Client) GetEmailLanguage(ctx context.Context) (result string, err error) {
	err = c.Call(ctx, "private/get_email_language", nil, &result)
	return
}

func (c *Client) GetNewAnnouncements(
	ctx context.Context,
) (result []models.Announcement, err error) {
	err = c.Call(ctx, "private/get_new_announcements", nil, &result)
	return
}

func (c *Client) GetPosition(
	ctx context.Context,
	params *models.GetPositionParams,
) (result models.Position, err error) {
	err = c.Call(ctx, "private/get_position", params, &result)
	return
}

func (c *Client) GetPositions(
	ctx context.Context,
	params *models.GetPositionsParams,
) (result []models.Position, err error) {
	err = c.Call(ctx, "private/get_positions", params, &result)
	return
}

func (c *Client) GetSubaccounts(
	ctx context.Context,
	params *models.GetSubaccountsParams,
) (result []models.Subaccount, err error) {
	err = c.Call(ctx, "private/get_subaccounts", params, &result)
	return
}

func (c *Client) SetAnnouncementAsRead(
	ctx context.Context,
	params *models.SetAnnouncementAsReadParams,
) (result string, err error) {
	err = c.Call(ctx, "private/set_announcement_as_read", params, &result)
	return
}

func (c *Client) SetEmailForSubaccount(
	ctx context.Context,
	params *models.SetEmailForSubaccountParams,
) (result string, err error) {
	err = c.Call(ctx, "private/set_email_for_subaccount", params, &result)
	return
}

func (c *Client) SetEmailLanguage(
	ctx context.Context,
	params *models.SetEmailLanguageParams,
) (result string, err error) {
	err = c.Call(ctx, "private/set_email_language", params, &result)
	return
}

func (c *Client) SetPasswordForSubaccount(
	ctx context.Context,
	params *models.SetPasswordForSubaccountParams,
) (result string, err error) {
	err = c.Call(ctx, "private/set_password_for_subaccount", params, &result)
	return
}

func (c *Client) ToggleNotificationsFromSubaccount(
	ctx context.Context,
	params *models.ToggleNotificationsFromSubaccountParams,
) (result string, err error) {
	err = c.Call(ctx, "private/toggle_notifications_from_subaccount", params, &result)
	return
}

func (c *Client) ToggleSubaccountLogin(
	ctx context.Context,
	params *models.ToggleSubaccountLoginParams,
) (result string, err error) {
	err = c.Call(ctx, "private/toggle_subaccount_login", params, &result)
	return
}
