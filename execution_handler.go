package main

type ExchangeHandler struct {
	market            MarketType
	balance           float64
	feePercentage     float64 //Fee percentage applied to each trade. 0.01 = 1%
	slipagePercentage float64 //Slipage percentage applied to each trade after execution
	percentageTraded  float64 //Percentage (0.01 = 1%) of the balance used to trade each individual single position.
	openPosition      *Position
	tradeHistory      []*Position
}

type MarketType int

const (
	CoinMarginedFutures MarketType = iota
	Futures
	Spot
)

func (handler *ExchangeHandler) CloseCurrentPosition() {
	//TODO:Closes the current open position
}

func (handler *ExchangeHandler) OpenNewPosition(tradeDirection Direction, entryPrice float64, leverage uint) {
	if tradeDirection == HOLD || entryPrice <= 0 {
		return
	}
}

func (handler *ExchangeHandler) ProcessLatestPrice(latestPrice *DataPoint) {

}
