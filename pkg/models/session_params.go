package models

type SetHeartbeatParams struct {
	Interval uint64 `json:"interval"`
}

type SessionParams struct {
	Scope string `json:"scope"`
}
