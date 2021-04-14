package pkg

type Backtester struct {
	eventQueue      EventQueue
	myStrategy      Strategy
	exchangeHandler *ExchangeHandler
	dataHandler     *DataHandler
}

type BacktestOptions struct {
	PriceWindow        uint
	Market             MarketType
	MakerFeePercentage float64
	TakerFeePercentage float64
	percentagePerTrade float64
}

//Event represents a action that will be processed by the eventloop
type Event interface{}

func NewBacktester(mystrategy Strategy, dataHandler *DataHandler) *Backtester {
	return &Backtester{
		exchangeHandler: NewExchangeHandler(USDFutures, 0.02, 0.04, 1),
		dataHandler:     dataHandler,
		myStrategy:      mystrategy,
	}
}

func NewCustomizedBacktester(mystrategy Strategy, dataHandler *DataHandler, options BacktestOptions) *Backtester {
	exchangeHandler := NewExchangeHandler(options.Market, options.MakerFeePercentage, options.TakerFeePercentage,
		options.percentagePerTrade)
	return &Backtester{
		exchangeHandler: exchangeHandler,
		dataHandler:     dataHandler,
		myStrategy:      mystrategy,
	}
}

func (bt *Backtester) Run() *Statistics {
	initialBalance := bt.exchangeHandler.balance
	var datapoints = bt.dataHandler.nextValues()

	for processed := bt.processNextEvent(); processed || datapoints != nil; processed = bt.processNextEvent() {
		if bt.eventQueue.IsEmpty() {
			if datapoints = bt.dataHandler.nextValues(); datapoints != nil {
				bt.eventQueue.AddEvent(datapoints)
			}
		}
	}
	return calculateStatistics(initialBalance, bt.exchangeHandler.tradeHistory)
}

//processNextEvent process the next event in the queue if the queue is not empty.
//returns a bool indicating if the event was processed or not
func (bt *Backtester) processNextEvent() bool {
	if !bt.eventQueue.HasNext() {
		return false
	}

	switch event := bt.eventQueue.NextEvent().(type) {
	case *AggregatedDataPoints:
		bt.processNewPriceEvt(event)
	case *OpenPositionEvt:
		bt.exchangeHandler.OpenMarketOrder(event.Direction, event.Leverage)
	case *StoplossEvt:
		bt.exchangeHandler.SetStoploss(event.Price)
	case *TakeProfitEvt:
		bt.exchangeHandler.SetTakeProfit(event.Price)

	}
	return true
}

func (bt *Backtester) processNewPriceEvt(newPrice *AggregatedDataPoints) {
	latestPrice := newPrice.datapoints[len(newPrice.datapoints)-1]
	bt.exchangeHandler.onPriceChange(latestPrice)

	if evt := bt.myStrategy.ProcessNextPriceData(newPrice.datapoints); evt != nil {
		bt.eventQueue.AddEvent(evt)
	}

	if bt.exchangeHandler.openPosition == nil {
		return
	}

	if evt := bt.myStrategy.SetStoploss(*bt.exchangeHandler.openPosition); evt != nil {
		bt.eventQueue.AddEvent(evt)
	}

	if evt := bt.myStrategy.SetTakeProfit(*bt.exchangeHandler.openPosition); evt != nil {
		bt.eventQueue.AddEvent(evt)
	}
}

func (bt *Backtester) SetSlippagePercentage(slippagePercent float64) {
	bt.exchangeHandler.SetSlipage(slippagePercent)
}
