package main

type Backtester struct {
	eventQueue      EventQueue
	myStrategy      Strategy
	exchangeHandler *ExchangeHandler
	dataHandler     *DataHandler
}

type BacktestOptions struct {
	PriceWindow        uint
	Market             MarketType
	Slipage            float64
	MakerFeePercentage float64
	TakerFeePercentage float64
	percentagePerTrade float64
}

//Event represents a action that will be processed by the eventloop
type Event interface{}

func NewBacktester(mystrategy Strategy, options BacktestOptions) *Backtester {
	exchangeHandler := NewExchangeHandler(options.Market, options.MakerFeePercentage, options.TakerFeePercentage,
		options.percentagePerTrade)
	return &Backtester{
		exchangeHandler: exchangeHandler,
		myStrategy:      mystrategy,
	}
}

func (bt *Backtester) Run() *Statistics {
	var datapoints = bt.dataHandler.nextValues()

	for processed := bt.processNextEvent(); !processed && datapoints == nil; {
		if datapoints = bt.dataHandler.nextValues(); datapoints != nil {
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

	switch event := bt.eventQueue.NextEvent().(type) {
	case DataPoint:
		bt.exchangeHandler.onPriceChange(event.Close)
	case OpenPositionEvt:
		bt.exchangeHandler.OpenMarketOrder(event.direction, event.leverage)
	case StoplossEvt:
		bt.exchangeHandler.SetStoploss(event.price)
	case TakeProfitEvt:
		bt.exchangeHandler.SetTakeProfit(event.price)
	}

	return true
}
