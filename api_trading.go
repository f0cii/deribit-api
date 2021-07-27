package deribit

import (
	"context"

	"github.com/KyberNetwork/deribit-api/models"
)

func (c *Client) Buy(ctx context.Context, params *models.BuyParams) (result models.BuyResponse, err error) {
	err = c.Call(ctx, "private/buy", params, &result)
	return
}

func (c *Client) Sell(ctx context.Context, params *models.SellParams) (result models.SellResponse, err error) {
	err = c.Call(ctx, "private/sell", params, &result)
	return
}

func (c *Client) Edit(ctx context.Context, params *models.EditParams) (result models.EditResponse, err error) {
	err = c.Call(ctx, "private/edit", params, &result)
	return
}

func (c *Client) Cancel(ctx context.Context, params *models.CancelParams) (result models.Order, err error) {
	err = c.Call(ctx, "private/cancel", params, &result)
	return
}

func (c *Client) CancelAll(ctx context.Context) (result string, err error) {
	err = c.Call(ctx, "private/cancel_all", nil, &result)
	return
}

func (c *Client) CancelAllByCurrency(ctx context.Context, params *models.CancelAllByCurrencyParams) (result string, err error) {
	err = c.Call(ctx, "private/cancel_all_by_currency", params, &result)
	return
}

func (c *Client) CancelAllByInstrument(ctx context.Context, params *models.CancelAllByInstrumentParams) (result string, err error) {
	err = c.Call(ctx, "private/cancel_all_by_instrument", params, &result)
	return
}

func (c *Client) ClosePosition(ctx context.Context, params *models.ClosePositionParams) (result models.ClosePositionResponse, err error) {
	err = c.Call(ctx, "private/close_position", params, &result)
	return
}

func (c *Client) GetMargins(ctx context.Context, params *models.GetMarginsParams) (result models.GetMarginsResponse, err error) {
	err = c.Call(ctx, "private/get_margins", params, &result)
	return
}

func (c *Client) GetOpenOrdersByCurrency(ctx context.Context, params *models.GetOpenOrdersByCurrencyParams) (result []models.Order, err error) {
	err = c.Call(ctx, "private/get_open_orders_by_currency", params, &result)
	return
}

func (c *Client) GetOpenOrdersByInstrument(ctx context.Context, params *models.GetOpenOrdersByInstrumentParams) (result []models.Order, err error) {
	err = c.Call(ctx, "private/get_open_orders_by_instrument", params, &result)
	return
}

func (c *Client) GetOrderHistoryByCurrency(ctx context.Context, params *models.GetOrderHistoryByCurrencyParams) (result []models.Order, err error) {
	err = c.Call(ctx, "private/get_order_history_by_currency", params, &result)
	return
}

func (c *Client) GetOrderHistoryByInstrument(ctx context.Context, params *models.GetOrderHistoryByInstrumentParams) (result []models.Order, err error) {
	err = c.Call(ctx, "private/get_order_history_by_instrument", params, &result)
	return
}

func (c *Client) GetOrderMarginByIDs(ctx context.Context, params *models.GetOrderMarginByIDsParams) (result models.GetOrderMarginByIDsResponse, err error) {
	err = c.Call(ctx, "private/get_order_margin_by_ids", params, &result)
	return
}

func (c *Client) GetOrderState(ctx context.Context, params *models.GetOrderStateParams) (result models.Order, err error) {
	err = c.Call(ctx, "private/get_order_state", params, &result)
	return
}

func (c *Client) GetUserTradesByCurrency(ctx context.Context, params *models.GetUserTradesByCurrencyParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call(ctx, "private/get_user_trades_by_currency", params, &result)
	return
}

func (c *Client) GetUserTradesByCurrencyAndTime(ctx context.Context, params *models.GetUserTradesByCurrencyAndTimeParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call(ctx, "private/get_user_trades_by_currency_and_time", params, &result)
	return
}

func (c *Client) GetUserTradesByInstrument(ctx context.Context, params *models.GetUserTradesByInstrumentParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call(ctx, "private/get_user_trades_by_instrument", params, &result)
	return
}

func (c *Client) GetUserTradesByInstrumentAndTime(ctx context.Context, params *models.GetUserTradesByInstrumentAndTimeParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call(ctx, "private/get_user_trades_by_instrument_and_time", params, &result)
	return
}

func (c *Client) GetUserTradesByOrder(ctx context.Context, params *models.GetUserTradesByOrderParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call(ctx, "private/get_user_trades_by_order", params, &result)
	return
}

func (c *Client) GetSettlementHistoryByInstrument(ctx context.Context, params *models.GetSettlementHistoryByInstrumentParams) (result models.GetSettlementHistoryResponse, err error) {
	err = c.Call(ctx, "private/get_settlement_history_by_instrument", params, &result)
	return
}

func (c *Client) GetSettlementHistoryByCurrency(ctx context.Context, params *models.GetSettlementHistoryByCurrencyParams) (result models.GetSettlementHistoryResponse, err error) {
	err = c.Call(ctx, "private/get_settlement_history_by_currency", params, &result)
	return
}
