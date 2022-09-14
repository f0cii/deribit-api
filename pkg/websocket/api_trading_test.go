package websocket

import (
	"context"
	"testing"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestBuy(t *testing.T) {
	expect := models.BuyResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       1966056,
				TradeID:        "ETH-2696083",
				Timestamp:      1590483938456,
				TickDirection:  0,
				Price:          203.3,
				MarkPrice:      203.28,
				InstrumentName: "ETH-PERPETUAL",
				IndexPrice:     203.33,
				Direction:      "buy",
				Amount:         40,
			},
		},
		Order: models.Order{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "filled",
			OrderType:           "market",
			OrderID:             "ETH-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "ETH-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Buy(
		context.Background(),
		&models.BuyParams{
			InstrumentName: "ETH-PERPETUAL",
			Amount:         40,
			Type:           "market",
			Label:          "market0000234",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestSell(t *testing.T) {
	expect := models.SellResponse{
		Trades: []models.Trade{},
		Order: models.Order{
			Trigger:             "last_price",
			TimeInForce:         "good_til_cancelled",
			Price:               145.61,
			OrderType:           "stop_limit",
			OrderState:          "untriggered",
			OrderID:             "ETH-SLTS-28",
			MaxShow:             123,
			LastUpdateTimestamp: 1550659803407,
			InstrumentName:      "ETH-PERPETUAL",
			Direction:           "sell",
			CreationTimestamp:   1550659803407,
			API:                 true,
			Amount:              123,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Sell(
		context.Background(),
		&models.SellParams{
			InstrumentName: "ETH-PERPETUAL",
			Amount:         123,
			Type:           "stop_limit",
			Price:          float64Pointer(145.61),
			TriggerPrice:   float64Pointer(145),
			Trigger:        "last_price",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestEdit(t *testing.T) {
	expect := models.EditResponse{
		Trades: []models.Trade{},
		Order: models.Order{
			TimeInForce:         "good_til_cancelled",
			Price:               0.1448,
			OrderState:          "open",
			OrderType:           "limit",
			OrderID:             "438994",
			MaxShow:             4,
			LastUpdateTimestamp: 1550585797677,
			InstrumentName:      "BTC-22FEB19-3500-C",
			Implv:               222,
			Direction:           "buy",
			CreationTimestamp:   1550585741277,
			Amount:              4,
			Advanced:            "implv",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Edit(
		context.Background(),
		&models.EditParams{
			OrderID:  "438994",
			Amount:   4,
			Price:    float64Pointer(222),
			Advanced: "implv",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestEditByLabel(t *testing.T) {
	expect := models.EditResponse{
		Trades: []models.Trade{},
		Order: models.Order{
			TimeInForce:         "good_til_cancelled",
			Price:               50111.0,
			OrderState:          "open",
			OrderType:           "limit",
			OrderID:             "94166",
			MaxShow:             150,
			LastUpdateTimestamp: 1616155550773,
			InstrumentName:      "BTC-PERPETUAL",
			Direction:           "buy",
			CreationTimestamp:   1616155547764,
			Label:               "test-12345",
			API:                 true,
			Amount:              150,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Edit(
		context.Background(),
		&models.EditParams{
			Label:  "test-12345",
			Amount: 150,
			Price:  float64Pointer(50111),
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancel(t *testing.T) {
	expect := models.Order{
		TimeInForce:         "good_til_cancelled",
		ProfitLoss:          0.00022929,
		Price:               203.3,
		OrderState:          "cancelled",
		OrderType:           "market",
		OrderID:             "ETH-584849853",
		MaxShow:             40,
		LastUpdateTimestamp: 1590483938456,
		Label:               "market0000234",
		InstrumentName:      "ETH-PERPETUAL",
		FilledAmount:        40,
		Direction:           "buy",
		CreationTimestamp:   1590483938456,
		Commission:          0.00014757,
		AveragePrice:        203.3,
		API:                 true,
		Amount:              40,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Cancel(
		context.Background(),
		&models.CancelParams{
			OrderID: "ETH-584849853",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancelAll(t *testing.T) {
	expect := uint(1)
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelAll(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancelAllByCurrency(t *testing.T) {
	expect := uint(1)
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelAllByCurrency(
		context.Background(),
		&models.CancelAllByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancelAllByInstrument(t *testing.T) {
	expect := uint(1)
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelAllByInstrument(
		context.Background(),
		&models.CancelAllByInstrumentParams{
			InstrumentName: "BTC-PERPETUAL",
			Type:           "all",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestCancelAllByLabel(t *testing.T) {
	expect := uint(1)
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.CancelAllByLabel(
		context.Background(),
		&models.CancelAllByInstrumentParams{
			Label: "test-12345",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestClosePosition(t *testing.T) {
	expect := models.ClosePositionResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       1966056,
				TradeID:        "ETH-2696083",
				Timestamp:      1590483938456,
				TickDirection:  0,
				Price:          203.3,
				MarkPrice:      203.28,
				InstrumentName: "ETH-PERPETUAL",
				IndexPrice:     203.33,
				Direction:      "buy",
				Amount:         40,
			},
		},
		Order: models.Order{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "filled",
			OrderType:           "market",
			OrderID:             "ETH-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "ETH-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.ClosePosition(
		context.Background(),
		&models.ClosePositionParams{
			InstrumentName: "ETH-PERPETUAL",
			Type:           "limit",
			Price:          float64Pointer(145.17),
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetMargins(t *testing.T) {
	expect := models.GetMarginsResponse{
		Buy:      0.01681367,
		Sell:     0.01680479,
		MaxPrice: 42.0,
		MinPrice: 42.0,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetMargins(
		context.Background(),
		&models.GetMarginsParams{
			InstrumentName: "BTC-PERPETUAL",
			Amount:         10000,
			Price:          3725,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOpenOrdersByCurrency(t *testing.T) {
	expect := []models.Order{
		{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "open",
			OrderType:           "market",
			OrderID:             "BTC-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "BTC-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOpenOrdersByCurrency(
		context.Background(),
		&models.GetOpenOrdersByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOpenOrdersByInstrument(t *testing.T) {
	expect := []models.Order{
		{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "open",
			OrderType:           "market",
			OrderID:             "BTC-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "BTC-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOpenOrdersByCurrency(
		context.Background(),
		&models.GetOpenOrdersByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOrderHistoryByCurrency(t *testing.T) {
	expect := []models.Order{
		{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "filled",
			OrderType:           "market",
			OrderID:             "BTC-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "BTC-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOrderHistoryByCurrency(
		context.Background(),
		&models.GetOrderHistoryByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOrderHistoryByInstrument(t *testing.T) {
	expect := []models.Order{
		{
			TimeInForce:         "good_til_cancelled",
			ProfitLoss:          0.00022929,
			Price:               203.3,
			OrderState:          "filled",
			OrderType:           "market",
			OrderID:             "BTC-584849853",
			MaxShow:             40,
			LastUpdateTimestamp: 1590483938456,
			Label:               "market0000234",
			InstrumentName:      "BTC-PERPETUAL",
			FilledAmount:        40,
			Direction:           "buy",
			CreationTimestamp:   1590483938456,
			Commission:          0.00014757,
			AveragePrice:        203.3,
			API:                 true,
			Amount:              40,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOrderHistoryByInstrument(
		context.Background(),
		&models.GetOrderHistoryByInstrumentParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOrderMarginByIDs(t *testing.T) {
	expect := models.GetOrderMarginByIDsResponse{
		{
			OrderID:       "ETH-349280",
			InitialMargin: 0.00091156,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOrderMarginByIDs(
		context.Background(),
		&models.GetOrderMarginByIDsParams{
			IDs: []string{"ETH-349280"},
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOrderState(t *testing.T) {
	expect := models.Order{
		TimeInForce:         "good_til_cancelled",
		ProfitLoss:          0.00022929,
		Price:               203.3,
		OrderState:          "filled",
		OrderType:           "market",
		OrderID:             "ETH-584849853",
		MaxShow:             40,
		LastUpdateTimestamp: 1590483938456,
		Label:               "market0000234",
		InstrumentName:      "ETH-PERPETUAL",
		FilledAmount:        40,
		Direction:           "buy",
		CreationTimestamp:   1590483938456,
		Commission:          0.00014757,
		AveragePrice:        203.3,
		API:                 true,
		Amount:              40,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOrderState(
		context.Background(),
		&models.GetOrderStateParams{
			OrderID: "ETH-584849853",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetUserTradesByCurrency(t *testing.T) {
	expect := models.GetUserTradesResponse{
		Trades: []models.UserTrade{
			{
				UnderlyingPrice: 204.5,
				TradeSeq:        3,
				TradeID:         "ETH-2696060",
				Timestamp:       1590480363130,
				TickDirection:   2,
				State:           "filled",
				Price:           0.361,
				OrderType:       "limit",
				OrderID:         "ETH-584827850",
				MarkPrice:       0.364585,
				Liquidity:       "T",
				InstrumentName:  "ETH-29MAY20-130-C",
				IndexPrice:      203.72,
				FeeCurrency:     "ETH",
				Fee:             0.002,
				Direction:       "sell",
				Amount:          5,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetUserTradesByCurrency(
		context.Background(),
		&models.GetUserTradesByCurrencyParams{
			Currency: "ETH",
			StartID:  "ETH-34066",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetUserTradesByCurrencyAndTime(t *testing.T) {
	expect := models.GetUserTradesResponse{
		Trades: []models.UserTrade{
			{
				UnderlyingPrice: 204.5,
				TradeSeq:        3,
				TradeID:         "ETH-2696060",
				Timestamp:       1590480363130,
				TickDirection:   2,
				State:           "filled",
				Price:           0.361,
				OrderType:       "limit",
				OrderID:         "ETH-584827850",
				MarkPrice:       0.364585,
				Liquidity:       "T",
				InstrumentName:  "ETH-29MAY20-130-C",
				IndexPrice:      203.72,
				FeeCurrency:     "ETH",
				Fee:             0.002,
				Direction:       "sell",
				Amount:          5,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetUserTradesByCurrencyAndTime(
		context.Background(),
		&models.GetUserTradesByCurrencyAndTimeParams{
			Currency:       "ETH",
			StartTimestamp: 1590480630731,
			EndTimestamp:   1510480630731,
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetUserTradesByInstrument(t *testing.T) {
	expect := models.GetUserTradesResponse{
		Trades: []models.UserTrade{
			{
				UnderlyingPrice: 204.5,
				TradeSeq:        1966042,
				TradeID:         "ETH-2696068",
				Timestamp:       1590480712800,
				TickDirection:   2,
				State:           "filled",
				Price:           203.8,
				OrderType:       "limit",
				OrderID:         "ETH-584827850",
				MarkPrice:       203.78,
				Liquidity:       "T",
				InstrumentName:  "ETH-PERPETUAL",
				IndexPrice:      203.72,
				FeeCurrency:     "ETH",
				Fee:             0.002,
				Direction:       "sell",
				Amount:          5,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetUserTradesByInstrument(
		context.Background(),
		&models.GetUserTradesByInstrumentParams{
			InstrumentName: "ETH-PERPETUAL",
			StartSeq:       1966042,
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetUserTradesByInstrumentAndTime(t *testing.T) {
	expect := models.GetUserTradesResponse{
		Trades: []models.UserTrade{
			{
				UnderlyingPrice: 204.5,
				TradeSeq:        1966042,
				TradeID:         "ETH-2696068",
				Timestamp:       1590480712800,
				TickDirection:   2,
				State:           "filled",
				Price:           203.8,
				OrderType:       "limit",
				OrderID:         "ETH-584827850",
				MarkPrice:       203.78,
				Liquidity:       "T",
				InstrumentName:  "ETH-PERPETUAL",
				IndexPrice:      203.72,
				FeeCurrency:     "ETH",
				Fee:             0.002,
				Direction:       "sell",
				Amount:          5,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetUserTradesByInstrumentAndTime(
		context.Background(),
		&models.GetUserTradesByInstrumentAndTimeParams{
			InstrumentName: "ETH-PERPETUAL",
			StartTimestamp: 1590470872894,
			EndTimestamp:   1590480872894,
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetUserTradesByOrder(t *testing.T) {
	expect := models.GetUserTradesResponse{
		Trades: []models.UserTrade{
			{
				UnderlyingPrice: 204.5,
				TradeSeq:        3,
				TradeID:         "ETH-2696060",
				Timestamp:       1590480363130,
				TickDirection:   2,
				State:           "filled",
				Price:           0.361,
				OrderType:       "limit",
				OrderID:         "ETH-584827850",
				MarkPrice:       0.364585,
				Liquidity:       "T",
				InstrumentName:  "ETH-29MAY20-130-C",
				IndexPrice:      203.72,
				FeeCurrency:     "ETH",
				Fee:             0.002,
				Direction:       "sell",
				Amount:          5,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetUserTradesByOrder(
		context.Background(),
		&models.GetUserTradesByOrderParams{
			OrderID: "ETH-584827850",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetSettlementHistoryByInstrument(t *testing.T) {
	expect := models.GetSettlementHistoryResponse{
		Settlements: []models.Settlement{
			{
				Type:              "settlement",
				Timestamp:         1550475692526,
				SessionProfitLoss: 0.038358299,
				ProfitLoss:        -0.001783937,
				Position:          -66,
				MarkPrice:         121.67,
				InstrumentName:    "ETH-22FEB19",
				IndexPrice:        119.8,
			},
		},
		Continuation: "xY7T6cusbMBNpH9SNmKb94jXSBxUPojJEdCPL4YociHBUgAhWQvEP",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetSettlementHistoryByInstrument(
		context.Background(),
		&models.GetSettlementHistoryByInstrumentParams{
			InstrumentName: "ETH-22FEB19",
			Type:           "settlement",
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetSettlementHistoryByCurrency(t *testing.T) {
	expect := models.GetSettlementHistoryResponse{
		Settlements: []models.Settlement{
			{
				Type:              "settlement",
				Timestamp:         1550475692526,
				SessionProfitLoss: 0.038358299,
				ProfitLoss:        -0.001783937,
				Position:          -66,
				MarkPrice:         121.67,
				InstrumentName:    "ETH-22FEB19",
				IndexPrice:        119.8,
			},
		},
		Continuation: "xY7T6cusbMBNpH9SNmKb94jXSBxUPojJEdCPL4YociHBUgAhWQvEP",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetSettlementHistoryByCurrency(
		context.Background(),
		&models.GetSettlementHistoryByCurrencyParams{
			Currency: "ETH",
			Type:     "settlement",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
