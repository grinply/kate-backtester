package kate

import "math"

//CoinMarket is the handler that allows trading simulation of COIN Margined crypto assets
type CoinMarket struct {
	Market   MarketType
	MakerFee float64
	TakerFee float64
}

func (marketHandler *CoinMarket) createPosition(tradeDirection Direction, currentPrice, balance, amountPerTrade float64, leverage uint) *Position {
	return nil
}

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000404726-P-L-calculations-Inverse-Contracts-#h_92ba55e9-4bbc-4879-a354-bc62eaa57d4d
func (marketHandler *CoinMarket) unrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * ((1 / position.EntryPrice) - (1 / lastTradedPrice))
	}
	return float64(position.Size) * ((1 / lastTradedPrice) - (1 / position.EntryPrice))
}

//CoinMarginedLiquidationPrice calculates the liquidation price for the position when trading COIN margined assets
//This calculation assumes the isolated trading position mode
//More info on https://help.bybit.com/hc/en-us/articles/360039261334-How-to-calculate-Liquidation-Price-Inverse-Contract
func (marketHandler *CoinMarket) liquidationPrice(position *Position) float64 {
	leverage := math.Max(1.0, float64(position.Leverage))
	if position.Direction == LONG {
		return (position.EntryPrice * leverage) / (leverage + 1 - (MMR * leverage))
	}
	return (position.EntryPrice * leverage) / (leverage - 1 + (MMR * leverage))
}

//marketFee calculates the fee applyed on market orders
func (marketHandler *CoinMarket) marketFee(position *Position) float64 {
	return position.Size * marketHandler.TakerFee
}

//limitFee calculates the fee applyed on limit orders
func (marketHandler *CoinMarket) limitFee(position *Position) float64 {
	return position.Size * marketHandler.TakerFee
}

func (marketHandler *CoinMarket) liquidationFee(position *Position) float64 {
	return 2 * marketHandler.marketFee(position)
}
