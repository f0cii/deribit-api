package deribit

import (
	"github.com/sumorf/deribit-api/models"
)

func (c *Client) Buy(params *models.BuyParams) (result models.BuyResponse, err error) {
	err = c.Call("private/buy", params, &result)
	return
}

func (c *Client) Sell(params *models.SellParams) (result models.SellResponse, err error) {
	err = c.Call("private/sell", params, &result)
	return
}

func (c *Client) Edit(params *models.EditParams) (result models.EditResponse, err error) {
	err = c.Call("private/edit", params, &result)
	return
}

func (c *Client) Cancel(params *models.CancelParams) (result models.Order, err error) {
	err = c.Call("private/cancel", params, &result)
	return
}

func (c *Client) CancelAll() (result string, err error) {
	err = c.Call("private/cancel_all", nil, &result)
	return
}

func (c *Client) CancelAllByCurrency(params *models.CancelAllByCurrencyParams) (result string, err error) {
	err = c.Call("private/cancel_all_by_currency", params, &result)
	return
}

func (c *Client) CancelAllByInstrument(params *models.CancelAllByInstrumentParams) (result string, err error) {
	err = c.Call("private/cancel_all_by_instrument", params, &result)
	return
}

func (c *Client) ClosePosition(params *models.ClosePositionParams) (result models.ClosePositionResponse, err error) {
	err = c.Call("private/close_position", params, &result)
	return
}

func (c *Client) GetMargins(params *models.GetMarginsParams) (result models.GetMarginsResponse, err error) {
	err = c.Call("private/get_margins", params, &result)
	return
}

func (c *Client) GetOpenOrdersByCurrency(params *models.GetOpenOrdersByCurrencyParams) (result []models.Order, err error) {
	err = c.Call("private/get_open_orders_by_currency", params, &result)
	return
}

func (c *Client) GetOpenOrdersByInstrument(params *models.GetOpenOrdersByInstrumentParams) (result []models.Order, err error) {
	err = c.Call("private/get_open_orders_by_instrument", params, &result)
	return
}

func (c *Client) GetOrderHistoryByCurrency(params *models.GetOrderHistoryByCurrencyParams) (result []models.Order, err error) {
	err = c.Call("private/get_order_history_by_currency", params, &result)
	return
}

func (c *Client) GetOrderHistoryByInstrument(params *models.GetOrderHistoryByInstrumentParams) (result []models.Order, err error) {
	err = c.Call("private/get_order_history_by_instrument", params, &result)
	return
}

func (c *Client) GetOrderMarginByIds(params *models.GetOrderMarginByIdsParams) (result models.GetOrderMarginByIdsResponse, err error) {
	err = c.Call("private/get_order_margin_by_ids", params, &result)
	return
}

func (c *Client) GetOrderState(params *models.GetOrderStateParams) (result models.Order, err error) {
	err = c.Call("private/get_order_state", params, &result)
	return
}

func (c *Client) GetStopOrderHistory(params *models.GetStopOrderHistoryParams) (result []models.StopOrder, err error) {
	err = c.Call("private/get_stop_order_history", params, &result)
	return
}

func (c *Client) GetUserTradesByCurrency(params *models.GetUserTradesByCurrencyParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call("private/get_user_trades_by_currency", params, &result)
	return
}

func (c *Client) GetUserTradesByCurrencyAndTime(params *models.GetUserTradesByCurrencyAndTimeParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call("private/get_user_trades_by_currency_and_time", params, &result)
	return
}

func (c *Client) GetUserTradesByInstrument(params *models.GetUserTradesByInstrumentParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call("private/get_user_trades_by_instrument", params, &result)
	return
}

func (c *Client) GetUserTradesByInstrumentAndTime(params *models.GetUserTradesByInstrumentAndTimeParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call("private/get_user_trades_by_instrument_and_time", params, &result)
	return
}

func (c *Client) GetUserTradesByOrder(params *models.GetUserTradesByOrderParams) (result models.GetUserTradesResponse, err error) {
	err = c.Call("private/get_user_trades_by_order", params, &result)
	return
}

func (c *Client) GetSettlementHistoryByInstrument(params *models.GetSettlementHistoryByInstrumentParams) (result models.GetSettlementHistoryResponse, err error) {
	err = c.Call("private/get_settlement_history_by_instrument", params, &result)
	return
}

func (c *Client) GetSettlementHistoryByCurrency(params *models.GetSettlementHistoryByCurrencyParams) (result models.GetSettlementHistoryResponse, err error) {
	err = c.Call("private/get_settlement_history_by_currency", params, &result)
	return
}
