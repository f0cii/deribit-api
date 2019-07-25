package models

type HelloParams struct {
	ClientName    string `json:"client_name"`
	ClientVersion string `json:"client_version"`
}
