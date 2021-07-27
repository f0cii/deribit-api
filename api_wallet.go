package deribit

import (
	"context"

	"github.com/KyberNetwork/deribit-api/models"
)

func (c *Client) CancelTransferByID(ctx context.Context, params *models.CancelTransferByIDParams) (result models.Transfer, err error) {
	err = c.Call(ctx, "private/cancel_transfer_by_id", params, &result)
	return
}

func (c *Client) CancelWithdrawal(ctx context.Context, params *models.CancelWithdrawalParams) (result models.Withdrawal, err error) {
	err = c.Call(ctx, "private/cancel_withdrawal", params, &result)
	return
}

func (c *Client) CreateDepositAddress(ctx context.Context, params *models.CreateDepositAddressParams) (result models.DepositAddress, err error) {
	err = c.Call(ctx, "private/create_deposit_address", params, &result)
	return
}

func (c *Client) GetCurrentDepositAddress(ctx context.Context, params *models.GetCurrentDepositAddressParams) (result models.DepositAddress, err error) {
	err = c.Call(ctx, "private/get_current_deposit_address", params, &result)
	return
}

func (c *Client) GetDeposits(ctx context.Context, params *models.GetDepositsParams) (result models.GetDepositsResponse, err error) {
	err = c.Call(ctx, "private/get_deposits", params, &result)
	return
}

func (c *Client) GetTransfers(ctx context.Context, params *models.GetTransfersParams) (result models.GetTransfersResponse, err error) {
	err = c.Call(ctx, "private/get_transfers", params, &result)
	return
}

func (c *Client) GetWithdrawals(ctx context.Context, params *models.GetWithdrawalsParams) (result []models.Withdrawal, err error) {
	err = c.Call(ctx, "private/get_withdrawals", params, &result)
	return
}

func (c *Client) Withdraw(ctx context.Context, params *models.WithdrawParams) (result models.Withdrawal, err error) {
	err = c.Call(ctx, "private/withdraw", params, &result)
	return
}
