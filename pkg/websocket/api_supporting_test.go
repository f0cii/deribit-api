package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetTime(t *testing.T) {
	expect := time.Now().UnixMilli()
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetTime(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestHello(t *testing.T) {
	expect := models.HelloResponse{Version: "1.0"}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Hello(
		context.Background(),
		&models.HelloParams{
			ClientName:    "test",
			ClientVersion: "1.0",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestTest(t *testing.T) {
	expect := models.TestResponse{Version: "1.0"}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Test(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
