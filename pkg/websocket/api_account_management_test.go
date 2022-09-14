package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAnnouncements(t *testing.T) {
	expect := []models.Announcement{
		{
			Title:           "test announcement",
			PublicationTime: uint64(time.Now().UnixMilli()),
			Important:       false,
			ID:              1,
			Body:            "this is a test announcement",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetAnnouncements(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestChangeSubaccountName(t *testing.T) {
	addResult(testClient.rpcConn, "test_account_name")
	res, err := testClient.ChangeSubaccountName(
		context.Background(),
		&models.ChangeSubaccountNameParams{
			Sid:  1,
			Name: "test_account_name",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, "test_account_name", res)
	}
}

func TestCreateSubaccount(t *testing.T) {
	expect := models.Subaccount{
		Email:                "test@example.test",
		ID:                   1,
		IsPassword:           false,
		LoginEnabled:         true,
		ReceiveNotifications: true,
		SystemName:           "Deribit test",
		TfaEnabled:           true,
		Type:                 "subaccount",
		Username:             "test",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CreateSubaccount(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestDisableTfaForSubaccount(t *testing.T) {
	addResult(testClient.rpcConn, successResponse)

	res, err := testClient.DisableTfaForSubaccount(
		context.Background(),
		&models.DisableTfaForSubaccountParams{
			Sid: 1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, successResponse, res)
	}
}

func TestGetAccountSummary(t *testing.T) {
	expect := models.AccountSummary{
		AvailableFunds:           1.0,
		AvailableWithdrawalFunds: 1.0,
		Balance:                  1.0,
		Currency:                 "BTC",
		DepositAddress:           "0x12345",
		Email:                    "test@example.com",
		Equity:                   1.0,
		ID:                       1,
		SystemName:               "Deribit test",
		TfaEnabled:               true,
		Type:                     "main",
		Username:                 "test",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetAccountSummary(
		context.Background(),
		&models.GetAccountSummaryParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetEmailLanguage(t *testing.T) {
	addResult(testClient.rpcConn, "test@example.com")

	res, err := testClient.GetEmailLanguage(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, "test@example.com", res)
	}
}

func TestGetNewAnnouncements(t *testing.T) {
	expect := []models.Announcement{
		{
			Title:           "test announcement",
			PublicationTime: uint64(time.Now().UnixMilli()),
			Important:       false,
			ID:              1,
			Body:            "this is a test announcement",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetNewAnnouncements(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetPosition(t *testing.T) {
	expect := models.Position{
		AveragePrice:              23000,
		AveragePriceUSD:           23000,
		Direction:                 "buy",
		EstimatedLiquidationPrice: 15000,
		FloatingProfitLoss:        1000,
		FloatingProfitLossUSD:     1000,
		IndexPrice:                24000,
		InitialMargin:             0.5,
		InstrumentName:            "BTC-PERPETUAL",
		Kind:                      "future",
		Leverage:                  10,
		MaintenanceMargin:         1,
		MarkPrice:                 24000,
		SettlementPrice:           23500,
		Size:                      24000,
		SizeCurrency:              1,
		TotalProfitLoss:           1000,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetPosition(
		context.Background(),
		&models.GetPositionParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetPositions(t *testing.T) {
	expect := []models.Position{
		{
			AveragePrice:              23000,
			AveragePriceUSD:           23000,
			Direction:                 "buy",
			EstimatedLiquidationPrice: 15000,
			FloatingProfitLoss:        1000,
			FloatingProfitLossUSD:     1000,
			IndexPrice:                24000,
			InitialMargin:             0.5,
			InstrumentName:            "BTC-PERPETUAL",
			Kind:                      "future",
			Leverage:                  10,
			MaintenanceMargin:         1,
			MarkPrice:                 24000,
			SettlementPrice:           23500,
			Size:                      24000,
			SizeCurrency:              1,
			TotalProfitLoss:           1000,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetPositions(
		context.Background(),
		&models.GetPositionsParams{
			Currency: "BTC",
			Kind:     "future",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetSubaccounts(t *testing.T) {
	expect := []models.Subaccount{
		{
			Email:                "test@example.test",
			ID:                   1,
			IsPassword:           false,
			LoginEnabled:         true,
			ReceiveNotifications: true,
			SystemName:           "Deribit test",
			TfaEnabled:           true,
			Type:                 "subaccount",
			Username:             "test",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetSubaccounts(
		context.Background(),
		&models.GetSubaccountsParams{},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestSetAnnouncementAsRead(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.SetAnnouncementAsRead(
		context.Background(),
		&models.SetAnnouncementAsReadParams{
			AnnouncementID: 1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestSetEmailForSubaccount(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.SetEmailForSubaccount(
		context.Background(),
		&models.SetEmailForSubaccountParams{
			Sid:   1,
			Email: "test@example.com",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestSetEmailLanguage(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.SetEmailLanguage(
		context.Background(),
		&models.SetEmailLanguageParams{
			Language: "english",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestSetPasswordForSubaccount(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.SetPasswordForSubaccount(
		context.Background(),
		&models.SetPasswordForSubaccountParams{
			Sid:      1,
			Password: "test@123456",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestToggleNotificationsFromSubaccount(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.ToggleNotificationsFromSubaccount(
		context.Background(),
		&models.ToggleNotificationsFromSubaccountParams{
			Sid:   1,
			State: true,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestToggleSubaccountLogin(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.ToggleSubaccountLogin(
		context.Background(),
		&models.ToggleSubaccountLoginParams{
			Sid:   1,
			State: "disable",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
