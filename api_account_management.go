package deribit

import "github.com/sumorf/deribit-api/models"

func (c *Client) GetPosition(params *models.GetPositionParams) (result models.Position, err error) {
	err = c.Call("private/get_position", params, &result)
	return
}

func (c *Client) GetPositions(params *models.GetPositionsParams) (result []models.Position, err error) {
	err = c.Call("private/get_positions", params, &result)
	return
}
