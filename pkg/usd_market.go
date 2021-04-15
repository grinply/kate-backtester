package pkg

import (
	"fmt"
	"math"
)

type USDMarket struct {
	Market   MarketType
	MakerFee float64
	TakerFee float64
}

func (marketHandler *USDMarket) createPosition(tradeDirection Direction, currentPrice, balance, amountPerTrade float64, leverage uint) *Position {
	usdMargin := balance * amountPerTrade

	newPosition := &Position{
		Direction:  tradeDirection,
		Margin:     usdMargin / currentPrice,
		Leverage:   leverage,
		EntryPrice: currentPrice,
	}
	newPosition.Size = math.Max(1.0, float64(leverage)) * newPosition.Margin
	newPosition.LiquidationPrice = marketHandler.liquidationPrice(newPosition)
	fmt.Println("Liquidation price: ", newPosition.LiquidationPrice)
	return newPosition
}

//liquidationPrice calculates the liquidation price for the positions when trading USD margined assets.
//This calculation assumes the isolated trading position mode
//More info on https://help.bybit.com/hc/en-us/articles/900000181046-Liquidation-Price-USDT-Contract
func (marketHandler *USDMarket) liquidationPrice(position *Position) float64 {
	leverage := math.Max(1.0, float64(position.Leverage))
	IMR := 1.0 / leverage
	if position.Direction == LONG {
		return position.EntryPrice * (1 - IMR + MMR)
	}
	return position.EntryPrice * (1 + IMR - MMR)
}

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000630066-P-L-calculations-USDT-Contract#Unrealized_P&L
func (marketHandler *USDMarket) unrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * (lastTradedPrice - position.EntryPrice)
	}
	return float64(position.Size) * (position.EntryPrice - lastTradedPrice)
}

//marketFee calculates the fee applyed on market orders
func (marketHandler *USDMarket) marketFee(position *Position) float64 {
	//If the close price is not zero it means the fee applyed is for closing the position
	//thus the price to apply the fee is the close price
	if position.ClosePrice > 0 {
		return position.Size * position.ClosePrice * marketHandler.TakerFee
	}
	return position.Size * position.EntryPrice * marketHandler.TakerFee
}

//limitFee calculates the fee applyed on limit orders
func (marketHandler *USDMarket) limitFee(position *Position) float64 {
	//If the close price is not zero it means the fee applyed is for closing the position
	//thus the price to apply the fee is the close price
	if position.ClosePrice > 0 {
		return position.Size * position.ClosePrice * marketHandler.MakerFee
	}
	return position.Size * position.EntryPrice * marketHandler.MakerFee
}

func (marketHandler *USDMarket) liquidationFee(position *Position) float64 {
	return 2 * marketHandler.marketFee(position)
}
