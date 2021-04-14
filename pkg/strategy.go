package pkg

type Direction int

const (
	LONG  Direction = iota //Open a new long position or close a short  if there is a position already opened
	SHORT                  //Opens a new short position or close a long if there is a position already opened
)

type Strategy interface {
	//ProcessNextPriceData check if it should open or close position
	//the latestPrices represent the most recent price data as defined in the window of the backtest
	ProcessNextPriceData(latestPrices []DataPoint) *OpenPositionEvt

	//SetStoploss [Optional] defines a price where a stoploss will be triggered closing the position in loss
	//The SetStoploss is called with updated unrealizedPnl everytime new price data is avaliable
	//A return value of -1 denotes that no stoploss will be set
	SetStoploss(openPosition Position) *StoplossEvt

	//SetTakeProfit defines a price where a takeprofit will be triggered closing the position in profit
	//The SetStoploss is called with updated unrealizedPnl everytime new price data is avaliable
	//A return value of -1 denotes that no takeprofit will be set
	SetTakeProfit(openPosition Position) *TakeProfitEvt
}
