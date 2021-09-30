module github.com/KyberNetwork/deribit-api

go 1.16

require (
	github.com/chuckpreslar/emission v0.0.0-20170206194824-a7ddd980baf9
	github.com/gorilla/websocket v1.4.1
	github.com/json-iterator/go v1.1.11
	github.com/shopspring/decimal v1.2.0
	github.com/sourcegraph/jsonrpc2 v0.1.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.18.1
)

replace github.com/sourcegraph/jsonrpc2 => github.com/KyberNetwork/jsonrpc2 v0.1.1-0.20210930035808-8a83c1f36cc0
