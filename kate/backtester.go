package kate

//Backtester allows backtesting trading strategies on crypto markets
type Backtester struct {
	eventQueue      EventQueue
	myStrategy      Strategy
	exchangeHandler *ExchangeHandler
	dataHandler     *DataHandler
}

//BacktestOptions is general settings for running a backtest
type BacktestOptions struct {
	TradedPair         string //Must follow the format: BTC/USD, ETH/USDT ...
	Market             MarketType
	MakerFeePercentage float64
	TakerFeePercentage float64
	percentagePerTrade float64
}

//Event represents a action that will be processed by the eventloop
type Event interface{}

//NewBacktester creates a new backtester instance that allows running trading simulations on crypto markets
func NewBacktester(mystrategy Strategy, dataHandler *DataHandler) *Backtester {
	return &Backtester{
		exchangeHandler: NewExchangeHandler(USDFutures, 0.02, 0.04, 1),
		dataHandler:     dataHandler,
		myStrategy:      mystrategy,
	}
}

//NewCustomizedBacktester creates a new customized backtester instance that allows running trading simulations on crypto markets
func NewCustomizedBacktester(mystrategy Strategy, dataHandler *DataHandler, options BacktestOptions) *Backtester {
	exchangeHandler := NewExchangeHandler(options.Market, options.MakerFeePercentage, options.TakerFeePercentage,
		options.percentagePerTrade)
	return &Backtester{
		exchangeHandler: exchangeHandler,
		dataHandler:     dataHandler,
		myStrategy:      mystrategy,
	}
}

//SetBalance defines the initial balance that will be used when trading
func (bt *Backtester) SetBalance(amount float64) {
	bt.exchangeHandler.balance = amount
}

//SetFixedTradeAmount defines a fixed value that will be used to open every position when trading
func (bt *Backtester) SetFixedTradeAmount(amount float64) {
	bt.exchangeHandler.fixedTradeAmount = amount
}

//Run executes a trading simulation for the provided configuration on the Backtester
func (bt *Backtester) Run() *Statistics {
	initialBalance := bt.exchangeHandler.balance

	for _, candle := range bt.dataHandler.Prices {
		for bt.eventQueue.HasNext() {
			bt.processNextEvent()
		}
		bt.eventQueue.AddEvent(candle)
	}

	return bt.calculateStatistics(initialBalance)
}

//processNextEvent process the next event in the queue if the queue is not empty.
func (bt *Backtester) processNextEvent() {
	switch event := bt.eventQueue.NextEvent().(type) {
	case DataPoint:
		bt.processNewPriceEvt(event)
	case *OpenPositionEvt:
		bt.exchangeHandler.OpenMarketOrder(event.Direction, event.Leverage)
	case *StoplossEvt:
		bt.exchangeHandler.SetStoploss(event.Price)
	case *TakeProfitEvt:
		bt.exchangeHandler.SetTakeProfit(event.Price)
	}
}

func (bt *Backtester) processNewPriceEvt(newPrice DataPoint) {
	bt.exchangeHandler.onPriceChange(newPrice)
	bt.myStrategy.PreProcessIndicators(newPrice)

	if bt.exchangeHandler.openPosition == nil {
		if evt := bt.myStrategy.OpenNewPosition(newPrice); evt != nil {
			bt.eventQueue.AddEvent(evt)
		}
	} else {
		if evt := bt.myStrategy.SetStoploss(*bt.exchangeHandler.openPosition); evt != nil {
			bt.eventQueue.AddEvent(evt)
		}

		if evt := bt.myStrategy.SetTakeProfit(*bt.exchangeHandler.openPosition); evt != nil {
			bt.eventQueue.AddEvent(evt)
		}
	}
}

//SetSlippagePercentage define a slippage that tries to better emulate the real trading market
func (bt *Backtester) SetSlippagePercentage(slippagePercent float64) {
	bt.exchangeHandler.SetSlipage(slippagePercent)
}
