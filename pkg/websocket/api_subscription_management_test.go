package websocket

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubscribeUnSubscribe(t *testing.T) {
	channels := []string{"book.BTC-PERPETUAL.raw", "user.trades.any.BTC.100ms"}
	addResult(testClient.rpcConn, []string{"book.BTC-PERPETUAL.raw"})
	addResult(testClient.rpcConn, []string{"user.trades.any.BTC.100ms"})
	err := testClient.Subscribe(channels)
	require.NoError(t, err)
	require.Len(t, testClient.subscriptions, 2)

	addResult(testClient.rpcConn, []string{"book.BTC-PERPETUAL.raw"})
	addResult(testClient.rpcConn, []string{"user.trades.any.BTC.100ms"})
	err = testClient.UnSubscribe(channels)
	require.NoError(t, err)
	require.Len(t, testClient.subscriptions, 0)
}
