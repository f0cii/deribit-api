package fix

import "github.com/quickfixgo/quickfix"

const (
	tagCancelOnDisconnect quickfix.Tag = 9001
	tagDeribitTradeID     quickfix.Tag = 100009
	tagDeribitLabel       quickfix.Tag = 100010
	tagMarkPrice          quickfix.Tag = 100090
	tagDeribitLiquidation quickfix.Tag = 100091
)
