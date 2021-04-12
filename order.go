package main

//Order types denotes how/when a execution of a position is made on the exchange.
//To know more check: https://www.binance.com/en/support/articles/360033779452
type OrderType int

const (
	MARKET OrderType = iota
	LIMIT
	STOPLIMIT
	STOPMARKET
	TRAILINGSTOP
)

//Order is a request of execution based on certain conditions to the exchange
type Order struct {
	Event
	id, leverage         uint
	direction            Direction
	orderType            OrderType
	totalQty             float64 //total value of the position including leverage
	margin, entryPrice   float64
	stoploss, takeprofit float64
}
