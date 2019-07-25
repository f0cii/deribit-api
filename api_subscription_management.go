package deribit

import "github.com/sumorf/deribit-api/models"

func (c *Client) PublicSubscribe(params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call("public/subscribe", params, &result)
	return
}

func (c *Client) PrivateSubscribe(params *models.SubscribeParams) (result models.SubscribeResponse, err error) {
	err = c.Call("private/subscribe", params, &result)
	return
}
