package websocket

import (
	"context"
	"testing"
	"time"

	"github.com/KyberNetwork/deribit-api/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestGetBookSummaryByCurrency(t *testing.T) {
	expect := []models.BookSummary{
		{
			AskPrice:          float64Pointer(20000),
			BaseCurrency:      "BTC",
			BidPrice:          float64Pointer(19999),
			CreationTimestamp: uint64(time.Now().UnixMilli()),
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetBookSummaryByCurrency(
		context.Background(),
		&models.GetBookSummaryByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetBookSummaryByInstrument(t *testing.T) {
	expect := []models.BookSummary{
		{
			AskPrice:          float64Pointer(20000),
			BaseCurrency:      "BTC",
			BidPrice:          float64Pointer(19999),
			CreationTimestamp: uint64(time.Now().UnixMilli()),
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetBookSummaryByInstrument(
		context.Background(),
		&models.GetBookSummaryByInstrumentParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetContractSize(t *testing.T) {
	expect := models.GetContractSizeResponse{
		ContractSize: 10,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetContractSize(
		context.Background(),
		&models.GetContractSizeParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetCurrencies(t *testing.T) {
	expect := []models.Currency{
		{
			CoinType: "crypto",
			Currency: "BTC",
		},
		{
			CoinType: "fiat",
			Currency: "USD",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetCurrencies(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetFundingChartData(t *testing.T) {
	expect := models.GetFundingChartDataResponse{
		CurrentInterest: 0.1,
		Data: []models.FundingChartData{
			{
				IndexPrice: 20000,
				Interest8H: 0.001,
				Timestamp:  uint64(time.Now().UnixMilli()),
			},
		},
		Interest8H: 0.001,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetFundingChartData(
		context.Background(),
		&models.GetFundingChartDataParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetHistoricalVolatility(t *testing.T) {
	expect := models.GetHistoricalVolatilityResponse{
		{
			Timestamp: uint64(time.Now().UnixMilli()),
			Value:     0.4,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetHistoricalVolatility(
		context.Background(),
		&models.GetHistoricalVolatilityParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetIndex(t *testing.T) {
	expect := models.GetIndexResponse{
		BTC: 23000,
		ETH: 1500,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetIndex(
		context.Background(),
		&models.GetIndexParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetInstrument(t *testing.T) {
	expect := models.Instrument{
		TickSize:             0.5,
		TakerCommission:      0.0005,
		SettlementPeriod:     "perpetual",
		QuoteCurrency:        "USD",
		MinTradeAmount:       10,
		MakerCommission:      0,
		Kind:                 "future",
		IsActive:             true,
		InstrumentName:       "BTC-PERPETUAL",
		InstrumentID:         124972,
		ExpirationTimestamp:  32503708800000,
		CreationTimestamp:    1534167754000,
		ContractSize:         10,
		BlockTradeCommission: 0.00025,
		BaseCurrency:         "BTC",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetInstrument(
		context.Background(),
		&models.GetInstrumentParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetInstruments(t *testing.T) {
	expect := []models.Instrument{
		{
			TickSize:             0.5,
			TakerCommission:      0.0005,
			SettlementPeriod:     "perpetual",
			QuoteCurrency:        "USD",
			MinTradeAmount:       10,
			MakerCommission:      0,
			Kind:                 "future",
			IsActive:             true,
			InstrumentName:       "BTC-PERPETUAL",
			InstrumentID:         124972,
			ExpirationTimestamp:  32503708800000,
			CreationTimestamp:    1534167754000,
			ContractSize:         10,
			BlockTradeCommission: 0.00025,
			BaseCurrency:         "BTC",
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetInstruments(
		context.Background(),
		&models.GetInstrumentsParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastSettlementsByCurrency(t *testing.T) {
	expect := models.GetLastSettlementsResponse{
		Settlements: []models.Settlement{
			{
				Type:              "delivery",
				Timestamp:         1663056000008,
				SessionProfitLoss: 7447,
				ProfitLoss:        0,
				Position:          75,
				MarkPrice:         0.23811436,
				InstrumentName:    "BTC-13SEP22-17000-C",
				IndexPrice:        22313.06,
			},
		},
		Continuation: "",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastSettlementsByCurrency(
		context.Background(),
		&models.GetLastSettlementsByCurrencyParams{
			Currency: "BTC",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastSettlementsByInstrument(t *testing.T) {
	expect := models.GetLastSettlementsResponse{
		Settlements: []models.Settlement{
			{
				Type:              "delivery",
				Timestamp:         1663056000008,
				SessionProfitLoss: 7447,
				ProfitLoss:        0,
				Position:          75,
				MarkPrice:         0.23811436,
				InstrumentName:    "BTC-PERPETUAL",
				IndexPrice:        22313.06,
			},
		},
		Continuation: "",
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastSettlementsByInstrument(
		context.Background(),
		&models.GetLastSettlementsByInstrumentParams{
			InstrumentName: "BTC-PERETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastTradesByCurrency(t *testing.T) {
	expect := models.GetLastTradesResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       81799266,
				TradeID:        "119946881",
				Timestamp:      1663058976578,
				TickDirection:  0,
				Price:          22232,
				MarkPrice:      22222.66,
				InstrumentName: "BTC-PERPETUAL",
				IndexPrice:     22211.47,
				Direction:      "buy",
				Amount:         990,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastTradesByCurrency(
		context.Background(),
		&models.GetLastTradesByCurrencyParams{
			Currency: "BTC",
			Count:    1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastTradesByCurrencyAndTime(t *testing.T) {
	expect := models.GetLastTradesResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       81799266,
				TradeID:        "119946881",
				Timestamp:      1663058976578,
				TickDirection:  0,
				Price:          22232,
				MarkPrice:      22222.66,
				InstrumentName: "BTC-PERPETUAL",
				IndexPrice:     22211.47,
				Direction:      "buy",
				Amount:         990,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastTradesByCurrencyAndTime(
		context.Background(),
		&models.GetLastTradesByCurrencyAndTimeParams{
			Currency:       "BTC",
			StartTimestamp: 1663058975578,
			EndTimestamp:   1663058977578,
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastTradesByInstrument(t *testing.T) {
	expect := models.GetLastTradesResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       81799266,
				TradeID:        "119946881",
				Timestamp:      1663058976578,
				TickDirection:  0,
				Price:          22232,
				MarkPrice:      22222.66,
				InstrumentName: "BTC-PERPETUAL",
				IndexPrice:     22211.47,
				Direction:      "buy",
				Amount:         990,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastTradesByInstrument(
		context.Background(),
		&models.GetLastTradesByInstrumentParams{
			InstrumentName: "BTC-PERPETUAL",
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetLastTradesByInstrumentAndTime(t *testing.T) {
	expect := models.GetLastTradesResponse{
		Trades: []models.Trade{
			{
				TradeSeq:       81799266,
				TradeID:        "119946881",
				Timestamp:      1663058976578,
				TickDirection:  0,
				Price:          22232,
				MarkPrice:      22222.66,
				InstrumentName: "BTC-PERPETUAL",
				IndexPrice:     22211.47,
				Direction:      "buy",
				Amount:         990,
			},
		},
		HasMore: false,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetLastTradesByInstrumentAndTime(
		context.Background(),
		&models.GetLastTradesByInstrumentAndTimeParams{
			InstrumentName: "BTC-PERPETUAL",
			StartTimestamp: 1663058975578,
			EndTimestamp:   1663058977578,
			Count:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetOrderBook(t *testing.T) {
	expect := models.GetOrderBookResponse{
		Timestamp: 1663059512573,
		Stats: models.Stats{
			Volume:      5489,
			PriceChange: float64Pointer(0.8791),
			Low:         20716.81,
			High:        22920,
		},
		State:           "open",
		SettlementPrice: 22290.36,
		OpenInterest:    3094866540,
		MinPrice:        21927.21,
		MaxPrice:        22595.04,
		MarkPrice:       22261.16,
		LastPrice:       22262,
		InstrumentName:  "BTC-PERPETUAL",
		IndexPrice:      22241.61,
		Funding8H:       0.0000483,
		CurrentFunding:  0.00037898,
		ChangeID:        14236384795,
		Bids: [][]float64{
			{
				22256,
				390,
			},
		},
		BestBidPrice:  22256,
		BestBidAmount: 390,
		BestAskPrice:  22262,
		BestAskAmount: 5900,
		Asks: [][]float64{
			{
				22262,
				5900,
			},
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetOrderBook(
		context.Background(),
		&models.GetOrderBookParams{
			InstrumentName: "BTC-PERPETUAL",
			Depth:          1,
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetTradeVolumes(t *testing.T) {
	expect := models.GetTradeVolumesResponse{
		{
			PutsVolume:    11059,
			FuturesVolume: 11824.0595099,
			CurrencyPair:  "btc_usd",
			CallsVolume:   15529.3,
		},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetTradeVolumes(
		context.Background(),
		&models.GetTradeVolumesParams{},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestGetTradingViewChartData(t *testing.T) {
	expect := models.GetTradingviewChartDataResponse{
		Volume: []float64{4.47511563},
		Ticks:  []uint64{1663059600000},
		Status: "ok",
		Open:   []float64{22259.5},
		Low:    []float64{22254},
		High:   []float64{22294.5},
		Cost:   []float64{99720},
		Close:  []float64{22281},
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.GetTradingviewChartData(
		context.Background(),
		&models.GetTradingviewChartDataParams{
			InstrumentName: "BTC-PERPETUAL",
			StartTimestamp: 1663059975640,
			EndTimestamp:   1663059975640,
			Resolution:     "10",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}

func TestTicker(t *testing.T) {
	expect := models.TickerResponse{
		Timestamp: 1663060531887,
		Stats: models.Stats{
			Volume:      5447,
			PriceChange: float64Pointer(5.3373),
			Low:         20716.81,
			High:        22920,
		},
		State:           "open",
		SettlementPrice: 22290.36,
		OpenInterest:    3094895590,
		MinPrice:        21954.2,
		MaxPrice:        22622.85,
		MarkPrice:       22287.66,
		LastPrice:       22292,
		InstrumentName:  "BTC-PERPETUAL",
		IndexPrice:      22279.72,
		Funding8H:       0.00005825,
		CurrentFunding:  0,
		BestBidPrice:    22283.5,
		BestBidAmount:   30,
		BestAskPrice:    22291.5,
		BestAskAmount:   140,
	}
	addResult(testClient.rpcConn, &expect)

	res, err := testClient.Ticker(
		context.Background(),
		&models.TickerParams{
			InstrumentName: "BTC-PERPETUAL",
		},
	)
	if assert.NoError(t, err) {
		assert.Equal(t, expect, res)
	}
}
