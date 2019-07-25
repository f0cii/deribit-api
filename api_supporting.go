package deribit

func (c *Client) GetTime() (timestamp int64, err error) {
	var result int64
	err = c.Call("public/get_time", nil, &result)
	if err != nil {
		return
	}
	timestamp = result
	return
}

func (c *Client) Test() (err error) {
	var result = struct {
		Version string `json:"version"`
	}{}
	err = c.Call("public/test", nil, &result)
	return
}
