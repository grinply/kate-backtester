package pkg

//Order types denotes how/when a execution of a position is made on the exchange.
//To know more check: https://www.binance.com/en/support/articles/360033779452
type OrderType int

const (
	MARKET OrderType = iota
	LIMIT
)

//Order is a request of execution based on certain conditions to the exchange
type OpenPositionEvt struct {
	Event
	Direction Direction
	Leverage  uint
	OrderType OrderType
}

type StoplossEvt struct {
	Event
	Price float64
}

type TakeProfitEvt struct {
	Event
	Price float64
}
