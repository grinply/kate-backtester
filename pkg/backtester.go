package main

type Backtester struct {
	eventQueue      EventQueue
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

//Event represents a action that will be processed by the eventloop
type Event interface{}

func NewBacktester(mystrategy Strategy, options BacktestOptions) *Backtester {
	return nil
}

func (bt *Backtester) Run() *Statistics {
	//if bt.dataHandler.x
	var datapoints *AggregatedDataPoints

	for processed := bt.processNextEvent(); !processed && datapoints == nil; {
		if datapoints = bt.dataHandler.NextValues(); !processed && datapoints != nil {
			bt.eventQueue.AddEvent(datapoints)
		}
	}
	return bt.calculateStatistics()
}

func (bt *Backtester) calculateStatistics() *Statistics {
	//TODO: calculate the statistics after each execution
	return nil
}

//processNextEvent process the next event in the queue if the queue is not empty.
//returns a bool indicating if the event was processed or not
func (bt *Backtester) processNextEvent() bool {
	if !bt.eventQueue.HasNext() {
		return false
	}

	switch bt.eventQueue.NextEvent().(type) {
	case DataPoint:
		//Process new tick data from the data stream
	case Order:
		//Process a new request order to exchange
		//case Fill:
		//Processs the order being filled (matched) on the exchange
	}

	return false
}

func (bt *Backtester) processLatestPrice(latestPriceData []DataPoint) {
	//bt.exchangeHandler.ProcessLatestPrice(&latestPriceData[len(latestPriceData)-1])
	bt.myStrategy.ProcessNextPriceData(latestPriceData)
	bt.myStrategy.SetStoploss(*bt.exchangeHandler.openPosition)
	bt.myStrategy.SetTakeProfit(*bt.exchangeHandler.openPosition)
}
