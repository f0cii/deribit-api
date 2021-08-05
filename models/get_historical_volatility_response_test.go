package models

import (
	"encoding/json"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestHistoricalVolatilityUnmarshalFunc(t *testing.T) {
	var hv HistoricalVolatility
	data := []byte(`[123456789, 2.2]`)

	err := json.Unmarshal(data, &hv)
	require.NoError(t, err)

	require.Equal(t, HistoricalVolatility{
		Timestamp: 123456789,
		Value:     decimal.NewFromFloat(2.2),
	}, hv)
}
