package deribit

import "github.com/frankrap/deribit-api/models"

func (c *Client) GetAnnouncements() (result []models.Announcement, err error) {
	err = c.Call("public/get_announcements", nil, &result)
	return
}

func (c *Client) ChangeSubaccountName(params *models.ChangeSubaccountNameParams) (result string, err error) {
	err = c.Call("private/change_subaccount_name", params, &result)
	return
}

func (c *Client) CreateSubaccount() (result models.Subaccount, err error) {
	err = c.Call("private/create_subaccount", nil, &result)
	return
}

func (c *Client) DisableTfaForSubaccount(params *models.DisableTfaForSubaccountParams) (result string, err error) {
	err = c.Call("private/disable_tfa_for_subaccount", params, &result)
	return
}

func (c *Client) GetAccountSummary(params *models.GetAccountSummaryParams) (result models.AccountSummary, err error) {
	err = c.Call("private/get_account_summary", params, &result)
	return
}

func (c *Client) GetEmailLanguage() (result string, err error) {
	err = c.Call("private/get_email_language", nil, &result)
	return
}

func (c *Client) GetNewAnnouncements() (result []models.Announcement, err error) {
	err = c.Call("private/get_new_announcements", nil, &result)
	return
}

func (c *Client) GetPosition(params *models.GetPositionParams) (result models.Position, err error) {
	err = c.Call("private/get_position", params, &result)
	return
}

func (c *Client) GetPositions(params *models.GetPositionsParams) (result []models.Position, err error) {
	err = c.Call("private/get_positions", params, &result)
	return
}

func (c *Client) GetSubaccounts(params *models.GetSubaccountsParams) (result []models.Subaccount, err error) {
	err = c.Call("private/get_subaccounts", params, &result)
	return
}

func (c *Client) GetSubaccountsDetails(params *models.GetSubaccountsDetailsParams) (result []models.SubaccountsDetails, err error) {
	err = c.Call("private/get_subaccounts_details", params, &result)
	return
}

func (c *Client) SetAnnouncementAsRead(params *models.SetAnnouncementAsReadParams) (result string, err error) {
	err = c.Call("private/set_announcement_as_read", params, &result)
	return
}

func (c *Client) SetEmailForSubaccount(params *models.SetEmailForSubaccountParams) (result string, err error) {
	err = c.Call("private/set_email_for_subaccount", params, &result)
	return
}

func (c *Client) SetEmailLanguage(params *models.SetEmailLanguageParams) (result string, err error) {
	err = c.Call("private/set_email_language", params, &result)
	return
}

func (c *Client) SetPasswordForSubaccount(params *models.SetPasswordForSubaccountParams) (result string, err error) {
	err = c.Call("private/set_password_for_subaccount", params, &result)
	return
}

func (c *Client) ToggleNotificationsFromSubaccount(params *models.ToggleNotificationsFromSubaccountParams) (result string, err error) {
	err = c.Call("private/toggle_notifications_from_subaccount", params, &result)
	return
}

func (c *Client) ToggleSubaccountLogin(params *models.ToggleSubaccountLoginParams) (result string, err error) {
	err = c.Call("private/toggle_subaccount_login", params, &result)
	return
}
