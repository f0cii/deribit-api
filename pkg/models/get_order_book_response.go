package models

type GetOrderBookResponse struct {
	Timestamp       uint64      `json:"timestamp"`
	Stats           Stats       `json:"stats"`
	State           string      `json:"state"`
	SettlementPrice float64     `json:"settlement_price"`
	MinPrice        float64     `json:"min_price"`
	MaxPrice        float64     `json:"max_price"`
	MarkPrice       float64     `json:"mark_price"`
	MarkIV          float64     `json:"mark_iv"`
	LastPrice       float64     `json:"last_price"`
	InstrumentName  string      `json:"instrument_name"`
	IndexPrice      float64     `json:"index_price"`
	Funding8H       float64     `json:"funding_8h"`
	CurrentFunding  float64     `json:"current_funding"`
	ChangeID        uint64      `json:"change_id"`
	BidIV           float64     `json:"bid_iv"`
	Bids            [][]float64 `json:"bids"`
	BestBidPrice    float64     `json:"best_bid_price"`
	BestBidAmount   float64     `json:"best_bid_amount"`
	BestAskPrice    float64     `json:"best_ask_price"`
	BestAskAmount   float64     `json:"best_ask_amount"`
	AskIV           float64     `json:"ask_iv"`
	Asks            [][]float64 `json:"asks"`
	DeliveryPrice   float64     `json:"delivery_price"`
	Greeks          Greeks      `json:"greeks"`
	InterestRate    float64     `json:"interest_rate"`
	OpenInterest    float64     `json:"open_interest"`
	UnderlyingIndex string      `json:"underlying_index"`
	UnderlyingPrice float64     `json:"underlying_price"`
}

type Greeks struct {
	Delta float64 `json:"delta"`
	Gamma float64 `json:"gamma"`
	RHO   float64 `json:"rho"`
	Theta float64 `json:"theta"`
	Vega  float64 `json:"vega"`
}
