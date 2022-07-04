package websocket

import (
	"context"

	"github.com/KyberNetwork/deribit-api/pkg/models"
)

func (c *Client) GetBookSummaryByCurrency(ctx context.Context, params *models.GetBookSummaryByCurrencyParams) (result []models.BookSummary, err error) {
	err = c.Call(ctx, "public/get_book_summary_by_currency", params, &result)
	return
}

func (c *Client) GetBookSummaryByInstrument(ctx context.Context, params *models.GetBookSummaryByInstrumentParams) (result []models.BookSummary, err error) {
	err = c.Call(ctx, "public/get_book_summary_by_instrument", params, &result)
	return
}

func (c *Client) GetContractSize(ctx context.Context, params *models.GetContractSizeParams) (result models.GetContractSizeResponse, err error) {
	err = c.Call(ctx, "public/get_contract_size", params, &result)
	return
}

func (c *Client) GetCurrencies(ctx context.Context) (result []models.Currency, err error) {
	err = c.Call(ctx, "public/get_currencies", nil, &result)
	return
}

func (c *Client) GetFundingChartData(ctx context.Context, params *models.GetFundingChartDataParams) (result models.GetFundingChartDataResponse, err error) {
	err = c.Call(ctx, "public/get_funding_chart_data", params, &result)
	return
}

func (c *Client) GetHistoricalVolatility(ctx context.Context, params *models.GetHistoricalVolatilityParams) (result models.GetHistoricalVolatilityResponse, err error) {
	err = c.Call(ctx, "public/get_historical_volatility", params, &result)
	return
}

func (c *Client) GetIndex(ctx context.Context, params *models.GetIndexParams) (result models.GetIndexResponse, err error) {
	err = c.Call(ctx, "public/get_index", params, &result)
	return
}

func (c *Client) GetInstrument(ctx context.Context, params *models.GetInstrumentParams) (result models.Instrument, err error) {
	err = c.Call(ctx, "public/get_instrument", params, &result)
	return
}

func (c *Client) GetInstruments(ctx context.Context, params *models.GetInstrumentsParams) (result []models.Instrument, err error) {
	err = c.Call(ctx, "public/get_instruments", params, &result)
	return
}

func (c *Client) GetLastSettlementsByCurrency(ctx context.Context, params *models.GetLastSettlementsByCurrencyParams) (result models.GetLastSettlementsResponse, err error) {
	err = c.Call(ctx, "public/get_last_settlements_by_currency", params, &result)
	return
}

func (c *Client) GetLastSettlementsByInstrument(ctx context.Context, params *models.GetLastSettlementsByInstrumentParams) (result models.GetLastSettlementsResponse, err error) {
	err = c.Call(ctx, "public/get_last_settlements_by_instrument", params, &result)
	return
}

func (c *Client) GetLastTradesByCurrency(ctx context.Context, params *models.GetLastTradesByCurrencyParams) (result models.GetLastTradesResponse, err error) {
	err = c.Call(ctx, "public/get_last_trades_by_currency", params, &result)
	return
}

func (c *Client) GetLastTradesByCurrencyAndTime(ctx context.Context, params *models.GetLastTradesByCurrencyAndTimeParams) (result models.GetLastTradesResponse, err error) {
	err = c.Call(ctx, "public/get_last_trades_by_currency_and_time", params, &result)
	return
}

func (c *Client) GetLastTradesByInstrument(ctx context.Context, params *models.GetLastTradesByInstrumentParams) (result models.GetLastTradesResponse, err error) {
	err = c.Call(ctx, "public/get_last_trades_by_instrument", params, &result)
	return
}

func (c *Client) GetLastTradesByInstrumentAndTime(ctx context.Context, params *models.GetLastTradesByInstrumentAndTimeParams) (result models.GetLastTradesResponse, err error) {
	err = c.Call(ctx, "public/get_last_trades_by_instrument_and_time", params, &result)
	return
}

func (c *Client) GetOrderBook(ctx context.Context, params *models.GetOrderBookParams) (result models.GetOrderBookResponse, err error) {
	err = c.Call(ctx, "public/get_order_book", params, &result)
	return
}

func (c *Client) GetTradeVolumes(ctx context.Context, params *models.GetTradeVolumesParams) (result models.GetTradeVolumesResponse, err error) {
	err = c.Call(ctx, "public/get_trade_volumes", params, &result)
	return
}

func (c *Client) GetTradingviewChartData(ctx context.Context, params *models.GetTradingviewChartDataParams) (result models.GetTradingviewChartDataResponse, err error) {
	err = c.Call(ctx, "public/get_tradingview_chart_data", params, &result)
	return
}

func (c *Client) Ticker(ctx context.Context, params *models.TickerParams) (result models.TickerResponse, err error) {
	err = c.Call(ctx, "public/ticker", params, &result)
	return
}
