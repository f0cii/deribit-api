package models

type SetHeartbeatParams struct {
	Interval float64 `json:"interval"`
}

type SessionParams struct {
	Scope string `json:"scope"`
}
