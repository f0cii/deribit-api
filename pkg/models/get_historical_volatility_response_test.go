package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHistoricalVolatilityUnmarshalFunc(t *testing.T) {
	t.Parallel()

	var hist HistoricalVolatility
	data := []byte(`[123456789, 2.2]`)

	err := json.Unmarshal(data, &hist)
	require.NoError(t, err)

	require.Equal(t, HistoricalVolatility{
		Timestamp: 123456789,
		Value:     2.2,
	}, hist)
}
