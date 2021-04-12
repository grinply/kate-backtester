package main

import (
	"fmt"
	"math"
)

type ExchangeHandler struct {
	market            MarketType
	balance           float64
	makerFee          float64 //Fee applied to limit orders - percentage applied is defined as 0.01 = 1%
	takerFee          float64 //Fee applied to market orders - percentage applied is defined as 0.01 = 1%
	slipagePercentage float64 //Slipage percentage applied to each trade after execution
	amountPerTrade    float64 //Percentage (0.01 = 1%) of the balance used to trade each individual single position.
	openPosition      *Position
	tradeHistory      []*Position
	currentPrice      float64 //price used as reference for latest price data - used to check if inputs are valid
}

type MarketType int

const (
	CoinMarginedFutures MarketType = iota
	USDFutures
	Spot
)

//NewExchangeHandler creates a new exchange handler that emulates exchange functionality
func NewExchangeHandler(market MarketType, makerFeePercent, takerFeePercent, percentagePerTrade float64) *ExchangeHandler {
	return &ExchangeHandler{
		market:            market,
		balance:           1000,
		slipagePercentage: 0.002,
		makerFee:          makerFeePercent / 100,
		takerFee:          takerFeePercent / 100,
		amountPerTrade:    percentagePerTrade / 100,
	}
}

//SetBalance sets the balance that will be used to trade
func (handler *ExchangeHandler) SetBalance(amount float64) {
	handler.balance = amount
}

//SetSlipage defines the slipage in the price on all orders of type market
func (handler *ExchangeHandler) SetSlipage(slipagePercent float64) {
	handler.slipagePercentage = slipagePercent
}

//ExecuteMarketOrder opens a new position with a market order if there is no positions already opened
func (handler *ExchangeHandler) OpenMarketOrder(tradeDirection Direction, leverage uint) error {
	if handler.openPosition != nil && handler.openPosition.Direction == tradeDirection {
		return fmt.Errorf("it is not possible to increase a already open position")
	}

	if handler.balance <= 1 {
		return fmt.Errorf("no more balance to trade")
	}

	margin := handler.balance * handler.amountPerTrade
	handler.openPosition = &Position{
		Direction:  tradeDirection,
		Margin:     margin,
		Leverage:   leverage,
		Size:       math.Max(1.0, float64(leverage)) * margin,
		EntryPrice: handler.currentPrice,
	}
	handler.openPosition.TotalFeePaid = handler.marketFee()
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
func (handler *ExchangeHandler) onPriceChange(newPrice float64) {
	handler.currentPrice = newPrice
	if handler.openPosition == nil {
		return
	}
	handler.updateUnrealizedPNL(newPrice)
	if handler.checkCloseLongs(newPrice) || handler.checkCloseShorts(newPrice) ||
		handler.checkLiquidation(newPrice) {
		return //Position closed sucessfuly
	}
}

func (handler *ExchangeHandler) updateUnrealizedPNL(latestPrice float64) {
	switch handler.market {
	case CoinMarginedFutures:
		handler.openPosition.UnrealizedPNL = CoinMarginedUnrealizedPNL(handler.openPosition, latestPrice)
	case USDFutures:
		handler.openPosition.UnrealizedPNL = USDMarginedUnrealizedPNL(handler.openPosition, latestPrice)
	}
}

func (handler *ExchangeHandler) checkCloseShorts(newPrice float64) bool {
	if handler.openPosition.Direction != SHORT {
		return false
	}

	if newPrice <= handler.openPosition.TakeProfit {
		handler.closePosition(newPrice, handler.limitFee())
		return true
	}

	if newPrice >= handler.openPosition.Stoploss {
		handler.closePosition(newPrice, handler.marketFee())
		return true
	}
	return false
}

func (handler *ExchangeHandler) checkCloseLongs(newPrice float64) bool {
	if handler.openPosition.Direction != LONG {
		return false
	}

	if newPrice >= handler.openPosition.TakeProfit {
		handler.closePosition(newPrice, handler.limitFee())
		return true
	}

	if newPrice <= handler.openPosition.Stoploss {
		handler.closePosition(newPrice, handler.marketFee())
		return true
	}
	return false
}

func (handler *ExchangeHandler) closePosition(latestPrice float64, closingFee float64) {
	handler.openPosition.ClosePrice = latestPrice
	handler.openPosition.TotalFeePaid += closingFee
	handler.openPosition.RealizedPNL = handler.openPosition.UnrealizedPNL - handler.openPosition.TotalFeePaid
	handler.openPosition.UnrealizedPNL = 0

	handler.tradeHistory = append(handler.tradeHistory, handler.openPosition)
	handler.openPosition = nil
}

//checkLiquidation verifies if a open position should be liquidated
func (handler *ExchangeHandler) checkLiquidation(newPrice float64) bool {
	if handler.openPosition.UnrealizedPNL == -(handler.openPosition.Size * 0.9) {
		handler.closePosition(newPrice, handler.liquidationFee())
		return true
	}
	return false
}

//marketFee calculates the fee applyed on market orders
func (handler *ExchangeHandler) marketFee() float64 {
	return handler.openPosition.Size * handler.takerFee
}

//limitFee calculates the fee applyed on limit orders
func (handler *ExchangeHandler) limitFee() float64 {
	return handler.openPosition.Size * handler.makerFee
}

func (handler *ExchangeHandler) liquidationFee() float64 {
	return 2 * handler.marketFee()
}
