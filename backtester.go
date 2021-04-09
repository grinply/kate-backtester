package main

type Backtester struct {
	myStrategy      Strategy
	exchangeHandler *ExchangeHandler
	dataHandler     *DataHandler
	priceWindow     int
}

type BacktestOptions struct {
	priceWindow       uint
	market            MarketType
	Slipage, TradeFee float64
}

type Stratistics struct {
	ROI         float64
	SharpeRatio float64
	TotalTrades uint
	WinRate     float64 //Percentage of wins
	MaxDrawdown float64 //Percentage for the maximum drawdown after applying the strategy
}

type Event int

const (
	FILL Event = iota
)

func NewBacktester(mystrategy Strategy, options BacktestOptions) *Backtester {
	return nil
}

func (bt *Backtester) Run() *Stratistics {
	var eventList []Event
	var latestPriceData []DataPoint

	for priceTick := bt.dataHandler.NextValue(); priceTick != nil && len(latestPriceData) < bt.priceWindow; {
		latestPriceData = append(latestPriceData, *priceTick)
	}

	for priceTick := bt.dataHandler.NextValue(); priceTick != nil; {
		latestPriceData = append(latestPriceData[1:], *priceTick)

		switch event := eventList[0]; event {
		case FILL:
			//Effectivelly executes the order on the exchange handler
		default:
			bt.processLatestPrice(latestPriceData)
		}
	}

	return &Stratistics{}
}

func (bt *Backtester) processLatestPrice(latestPriceData []DataPoint) {
	bt.exchangeHandler.ProcessLatestPrice(&latestPriceData[len(latestPriceData)-1])
	bt.myStrategy.ProcessNextPriceData(latestPriceData)
	bt.myStrategy.SetStoploss(*bt.exchangeHandler.openPosition)
	bt.myStrategy.SetTakeProfit(*bt.exchangeHandler.openPosition)
}
