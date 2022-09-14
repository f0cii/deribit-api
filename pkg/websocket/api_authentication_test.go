package websocket

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthLogout(t *testing.T) {
	client := newClient()
	err := client.Start()
	require.NoError(t, err)

	err = client.Logout(context.Background())
	require.NoError(t, err)
}
