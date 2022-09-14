package websocket

import (
	"context"
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestCancelTransferByID(t *testing.T) {
	expect := models.Transfer{
		Amount:           1,
		CreatedTimestamp: 1550579457727,
		Currency:         "BTC",
		Direction:        "payment",
		ID:               1,
		State:            "cancelled",
		Type:             "user",
		UpdatedTimestamp: 1550579457727,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelTransferByID(
		context.Background(),
		&models.CancelTransferByIDParams{
			Currency: "BTC",
			ID:       1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancelWithdrawal(t *testing.T) {
	expect := models.Withdrawal{
		Address:          "2NBqqD5GRJ8wHy1PYyCXTe9ke5226FhavBz",
		Amount:           1,
		CreatedTimestamp: 1550571443070,
		Currency:         "BTC",
		Fee:              0.0001,
		ID:               1,
		Priority:         0.15,
		State:            "cancelled",
		UpdatedTimestamp: 1550571443070,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelWithdrawal(
		context.Background(),
		&models.CancelWithdrawalParams{
			Currency: "BTC",
			ID:       1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCreateDepositAddress(t *testing.T) {
	expect := models.DepositAddress{
		Address:           "2N8udZGBc1hLRCFsU9kGwMPpmYUwMFTuCwB",
		CreationTimestamp: 1550575165170,
		Currency:          "BTC",
		Type:              "deposit",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CreateDepositAddress(
		context.Background(),
		&models.CreateDepositAddressParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetCurrentDepositAddress(t *testing.T) {
	expect := models.DepositAddress{
		Address:           "2N8udZGBc1hLRCFsU9kGwMPpmYUwMFTuCwB",
		CreationTimestamp: 1550575165170,
		Currency:          "BTC",
		Type:              "deposit",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetCurrentDepositAddress(
		context.Background(),
		&models.GetCurrentDepositAddressParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetDeposits(t *testing.T) {
	expect := models.GetDepositsResponse{
		Count: 1,
		Data: []models.Deposit{
			{
				Address:           "2N35qDKDY22zmJq9eSyiAerMD4enJ1xx6ax",
				Amount:            5,
				Currency:          "BTC",
				ReceivedTimestamp: 1549295017670,
				State:             "completed",
				TransactionID:     "230669110fdaf0a0dbcdc079b6b8b43d5af29cc73683835b9bc6b3406c065fda",
				UpdatedTimestamp:  1549295130159,
			},
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetDeposits(
		context.Background(),
		&models.GetDepositsParams{
			Currency: "BTC",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetTransfers(t *testing.T) {
	expect := models.GetTransfersResponse{
		Count: 1,
		Data: []models.Transfer{
			{
				Amount:           1,
				CreatedTimestamp: 1550579457727,
				Currency:         "BTC",
				Direction:        "payment",
				ID:               2,
				OtherSide:        "2MzyQc5Tkik61kJbEpJV5D5H9VfWHZK9Sgy",
				State:            "prepared",
				Type:             "user",
				UpdatedTimestamp: 1550579457727,
			},
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetTransfers(
		context.Background(),
		&models.GetTransfersParams{
			Currency: "BTC",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetWithdrawals(t *testing.T) {
	expect := []models.Withdrawal{
		{
			Address:          "2NBqqD5GRJ8wHy1PYyCXTe9ke5226FhavBz",
			Amount:           1,
			CreatedTimestamp: 1550571443070,
			Currency:         "BTC",
			Fee:              0.0001,
			ID:               1,
			Priority:         0.15,
			State:            "unconfirmed",
			UpdatedTimestamp: 1550571443070,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetWithdrawals(
		context.Background(),
		&models.GetWithdrawalsParams{
			Currency: "BTC",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestWithdraw(t *testing.T) {
	expect := models.Withdrawal{
		Address:          "2NBqqD5GRJ8wHy1PYyCXTe9ke5226FhavBz",
		Amount:           1,
		CreatedTimestamp: 1550571443070,
		Currency:         "BTC",
		Fee:              0.0001,
		ID:               1,
		Priority:         0.15,
		State:            "unconfirmed",
		UpdatedTimestamp: 1550571443070,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Withdraw(
		context.Background(),
		&models.WithdrawParams{
			Currency: "BTC",
			Address:  "2NBqqD5GRJ8wHy1PYyCXTe9ke5226FhavBz",
			Amount:   1,
			Priority: "0.15",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
