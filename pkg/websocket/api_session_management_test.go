package websocket

import (
	"context"
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestSetHeartbeat(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.SetHeartbeat(context.Background(), &models.SetHeartbeatParams{})
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestDisableHeartbeat(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.DisableHeartbeat(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestEnableCancelOnDisconnect(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.EnableCancelOnDisconnect(context.Background(), &models.SessionParams{})
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestDisableCancelOnDisconnect(t *testing.T) {
	expect := successResponse
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.DisableCancelOnDisconnect(context.Background(), &models.SessionParams{})
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
