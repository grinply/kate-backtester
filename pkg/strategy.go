package pkg

//Direction denotes the trade direction that a Position can have
type Direction int

const (
	//LONG denotes the trade direction that earns when the price moves upwards
	LONG Direction = iota
	//SHORT denotes the trade direction that earns when the price moves downwards
	SHORT
)

//Strategy defines how/when trades should be opened and how stoploss/takeprofits should be set in a simulated run
type Strategy interface {
	//OpenNewPosition check if a position should be open
	//the latestPrices represent the most recent price data as defined in the window of the backtest
	OpenNewPosition(latestPrices []DataPoint) *OpenPositionEvt

	//SetStoploss [Optional] defines a price where a stoploss will be triggered closing the position in loss
	//The SetStoploss is called with updated unrealizedPnl everytime new price data is available
	//A return value of -1 denotes that no stoploss will be set
	SetStoploss(openPosition Position) *StoplossEvt

	//SetTakeProfit defines a price where a takeprofit will be triggered closing the position in profit
	//The SetStoploss is called with updated unrealizedPnl everytime new price data is available
	//A return value of -1 denotes that no takeprofit will be set
	SetTakeProfit(openPosition Position) *TakeProfitEvt
}
