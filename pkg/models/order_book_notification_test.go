package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalFunc(t *testing.T) {
	t.Parallel()

	var bookNotif OrderBookNotificationItem
	data := []byte(`["new", 1.1, 2.2]`)

	err := json.Unmarshal(data, &bookNotif)
	require.NoError(t, err)

	require.Equal(t, OrderBookNotificationItem{
		Action: "new",
		Price:  1.1,
		Amount: 2.2,
	}, bookNotif)
}
