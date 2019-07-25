package deribit

import "github.com/sumorf/deribit-api/models"

func (c *Client) GetTime() (result int64, err error) {
	err = c.Call("public/get_time", nil, &result)
	return
}

func (c *Client) Hello(params *models.HelloParams) (result models.HelloResponse, err error) {
	err = c.Call("public/hello", params, &result)
	return
}

func (c *Client) Test() (result models.TestResponse, err error) {
	err = c.Call("public/test", nil, &result)
	return
}
