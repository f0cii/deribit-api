package deribit

import "github.com/sumorf/deribit-api/models"

func (c *Client) CancelTransferById(params *models.CancelTransferByIdParams) (result models.Transfer, err error) {
	err = c.Call("private/cancel_transfer_by_id", params, &result)
	return
}

func (c *Client) CancelWithdrawal(params *models.CancelWithdrawalParams) (result models.Withdrawal, err error) {
	err = c.Call("private/cancel_withdrawal", params, &result)
	return
}

func (c *Client) CreateDepositAddress(params *models.CreateDepositAddressParams) (result models.DepositAddress, err error) {
	err = c.Call("private/create_deposit_address", params, &result)
	return
}

func (c *Client) GetCurrentDepositAddress(params *models.GetCurrentDepositAddressParams) (result models.DepositAddress, err error) {
	err = c.Call("private/get_current_deposit_address", params, &result)
	return
}

func (c *Client) GetDeposits(params *models.GetDepositsParams) (result models.GetDepositsResponse, err error) {
	err = c.Call("private/get_deposits", params, &result)
	return
}

func (c *Client) GetTransfers(params *models.GetTransfersParams) (result models.GetTransfersResponse, err error) {
	err = c.Call("private/get_transfers", params, &result)
	return
}

func (c *Client) GetWithdrawals(params *models.GetWithdrawalsParams) (result []models.Withdrawal, err error) {
	err = c.Call("private/get_withdrawals", params, &result)
	return
}

func (c *Client) Withdraw(params *models.WithdrawParams) (result models.Withdrawal, err error) {
	err = c.Call("private/withdraw", params, &result)
	return
}
