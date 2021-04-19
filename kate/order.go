package kate

//OrderType denotes how/when a execution of a position is made on the exchange.
//To know more check: https://www.binance.com/en/support/articles/360033779452
type OrderType int

const (
	//MARKET is a type of order that executes immediately at the best price
	MARKET OrderType = iota
	//LIMIT is a type of order that executes when the traget price is reached
	LIMIT
)

//OpenPositionEvt is a event to open a simulated position
type OpenPositionEvt struct {
	Event
	Direction Direction
	Leverage  uint
	OrderType OrderType
}

//StoplossEvt is a event to set a stoploss
type StoplossEvt struct {
	Event
	Price float64
}

//TakeProfitEvt is a event to set a takeprofit
type TakeProfitEvt struct {
	Event
	Price float64
}
