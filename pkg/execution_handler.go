package pkg

import (
	"fmt"
)

type ExchangeHandler struct {
	marketHandler  MarketHandler
	balance        float64
	makerFee       float64 //Fee applied to limit orders - percentage applied is defined as 0.01 = 1%
	takerFee       float64 //Fee applied to market orders - percentage applied is defined as 0.01 = 1%
	slippage       float64 //Slipage percentage applied to each trade after execution
	amountPerTrade float64 //Percentage (0.01 = 1%) of the balance used to trade each individual single position.
	openPosition   *Position
	tradeHistory   []*Position
	currentPrice   float64 //price used as reference for latest price data - used to check if inputs are valid
}

type MarketType int
type PositionTransition int

const (
	CoinMarginedFutures MarketType = iota
	USDFutures
	Spot
)

const (
	MakerTransition PositionTransition = iota
	TakerTransition
	Liquidation
)

//NewExchangeHandler creates a new exchange handler that emulates exchange functionality
func NewExchangeHandler(market MarketType, makerFeePercent, takerFeePercent, percentagePerTrade float64) *ExchangeHandler {
	handler := &ExchangeHandler{
		balance:        1000,
		slippage:       0.0002,
		makerFee:       makerFeePercent / 100,
		takerFee:       takerFeePercent / 100,
		amountPerTrade: percentagePerTrade / 100,
	}
	handler.marketHandler = newMarketHandler(market, handler.makerFee, handler.takerFee)
	return handler
}

//SetBalance sets the balance that will be used to trade
func (handler *ExchangeHandler) SetBalance(amount float64) {
	handler.balance = amount
}

//SetSlipage defines the slipage in the price on all orders of type market
func (handler *ExchangeHandler) SetSlipage(slipagePercent float64) {
	handler.slippage = slipagePercent / 100
	//TODO: implement logic to apply the slipagge when opening executing (open/close) market orders
}

//ExecuteMarketOrder opens a new position with a market order if there is no positions already opened
func (handler *ExchangeHandler) OpenMarketOrder(tradeDirection Direction, leverage uint) error {
	if handler.openPosition != nil {
		return fmt.Errorf("there is a position already opened")
	}

	if handler.balance <= 1 {
		return fmt.Errorf("no more balance to trade")
	}

	handler.openPosition = handler.marketHandler.createPosition(tradeDirection, handler.currentPrice,
		handler.balance, handler.amountPerTrade, leverage)
	handler.openPosition.TotalFeePaid = handler.fee(TakerTransition)
	return nil
}

func (handler *ExchangeHandler) OpenLimitOrder(tradeDirection Direction, entryPrice, margin float64, leverage uint) error {
	if handler.openPosition != nil && handler.openPosition.Direction == tradeDirection {
		return fmt.Errorf("it is not possible to increase a already open position")
	}
	return nil
}

//SetStoploss defines a stoploss that closes the open position completely when the price is reached.
//The stoploss triggered is a market order
func (handler *ExchangeHandler) SetStoploss(price float64) error {
	if handler.openPosition == nil {
		return fmt.Errorf("there is no positions open to set a stoploss")
	}

	if handler.openPosition.Direction == LONG && price > handler.currentPrice {
		return fmt.Errorf("the stoploss must be lower than the current price for long positions")
	}

	if handler.openPosition.Direction == SHORT && price < handler.currentPrice {
		return fmt.Errorf("the stoploss must be higher than the current price for short positions")
	}

	handler.openPosition.Stoploss = price
	return nil
}

func (handler *ExchangeHandler) SetTakeProfit(price float64) error {
	if handler.openPosition == nil {
		return fmt.Errorf("there is no positions open to set a takeprofit")
	}

	if handler.openPosition.Direction == LONG && price < handler.currentPrice {
		return fmt.Errorf("the takeprofit must be higher than the current price for long positions")
	}

	if handler.openPosition.Direction == SHORT && price > handler.currentPrice {
		return fmt.Errorf("the takeprofit must be lower than the current price for short positions")
	}

	handler.openPosition.TakeProfit = price
	return nil
}

//OnPriceChange emulates the price change for the asset.
//Positions may be closed by: take profit, stoploss or liquidations.
func (handler *ExchangeHandler) onPriceChange(newPrice DataPoint) {
	handler.currentPrice = newPrice.Close
	if handler.openPosition == nil {
		return
	}

	if handler.checkCloseLongs(newPrice) || handler.checkCloseShorts(newPrice) ||
		handler.checkLiquidation(newPrice) {
		return //Position closed sucessfuly
	}
	//handler.updateUnrealizedPNL(newPrice.Close)
}

func (handler *ExchangeHandler) updateUnrealizedPNL(latestPrice float64) {
	handler.openPosition.UnrealizedPNL = handler.marketHandler.unrealizedPNL(handler.openPosition, latestPrice)
}

func (handler *ExchangeHandler) checkCloseShorts(newPrice DataPoint) bool {
	if handler.openPosition.Direction != SHORT {
		return false
	}

	if handler.openPosition.TakeProfit > 0 && newPrice.Low <= handler.openPosition.TakeProfit {
		handler.closePosition(handler.openPosition.TakeProfit, MakerTransition)
		return true
	}

	if handler.openPosition.Stoploss > 0 && newPrice.High >= handler.openPosition.Stoploss {
		handler.closePosition(handler.openPosition.Stoploss, TakerTransition)
		return true
	}
	return false
}

func (handler *ExchangeHandler) checkCloseLongs(newPrice DataPoint) bool {
	if handler.openPosition.Direction != LONG {
		return false
	}

	if handler.openPosition.TakeProfit > 0 && newPrice.High >= handler.openPosition.TakeProfit {
		handler.closePosition(handler.openPosition.TakeProfit, MakerTransition)
		return true
	}

	if handler.openPosition.Stoploss > 0 && newPrice.Low <= handler.openPosition.Stoploss {
		handler.closePosition(handler.openPosition.Stoploss, TakerTransition)
		return true
	}
	return false
}

func (handler *ExchangeHandler) closePosition(closePrice float64, transition PositionTransition) {
	handler.updateUnrealizedPNL(closePrice)
	handler.openPosition.ClosePrice = closePrice
	handler.openPosition.TotalFeePaid += handler.fee(transition)

	handler.openPosition.RealizedPNL = handler.openPosition.UnrealizedPNL - handler.openPosition.TotalFeePaid
	handler.openPosition.UnrealizedPNL = 0
	handler.balance += handler.openPosition.RealizedPNL

	handler.tradeHistory = append(handler.tradeHistory, handler.openPosition)
	handler.openPosition = nil
}

//checkLiquidation verifies if a open position should be liquidated
func (handler *ExchangeHandler) checkLiquidation(newPrice DataPoint) bool {
	fmt.Printf("%f <= %f\n", handler.openPosition.LiquidationPrice, newPrice.High)

	if handler.openPosition.Direction == LONG && handler.openPosition.LiquidationPrice >= newPrice.Low {
		handler.closePosition(handler.openPosition.LiquidationPrice, Liquidation)
		return true
	}

	if handler.openPosition.Direction == SHORT && handler.openPosition.LiquidationPrice <= newPrice.High {
		fmt.Println("Liquidated short")
		handler.closePosition(handler.openPosition.LiquidationPrice, Liquidation)
	}
	return false
}

func (handler *ExchangeHandler) fee(transition PositionTransition) float64 {
	switch transition {
	case MakerTransition:
		return handler.marketHandler.limitFee(handler.openPosition)
	case TakerTransition:
		return handler.marketHandler.marketFee(handler.openPosition)
	case Liquidation:
		return handler.marketHandler.liquidationFee(handler.openPosition)
	}
	return 0.0
}
