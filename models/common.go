package models

// Direction direction, `buy` or `sell`
const (
	DirectionBuy  = "buy"
	DirectionSell = "sell"
)

// OrderState order state, `"open"`, `"filled"`, `"rejected"`, `"cancelled"`, `"untriggered"`
const (
	OrderStateOpen        = "open"
	OrderStateFilled      = "filled"
	OrderStateRejected    = "rejected"
	OrderStateCancelled   = "cancelled"
	OrderStateUntriggered = "untriggered"
)

// OrderType order type, `"limit"`, `"market"`, `"stop_limit"`, `"stop_market"`
const (
	OrderTypeLimit      = "limit"
	OrderTypeMarket     = "market"
	OrderTypeStopLimit  = "stop_limit"
	OrderTypeStopMarket = "stop_market"
)
