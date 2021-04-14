package pkg

import "math"

//Maintenance margin rate
const MMR = 0.005

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000404726-P-L-calculations-Inverse-Contracts-#h_92ba55e9-4bbc-4879-a354-bc62eaa57d4d
func CoinMarginedUnrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * ((1 / position.EntryPrice) - (1 / lastTradedPrice))
	}
	return float64(position.Size) * ((1 / lastTradedPrice) - (1 / position.EntryPrice))
}

//UnrealizedPNL calculates the unrealized profit or loss ( in absolute values ) for the provided position
//to know more about the pnl calculation see: https://help.bybit.com/hc/en-us/articles/900000630066-P-L-calculations-USDT-Contract#Unrealized_P&L
func USDMarginedUnrealizedPNL(position *Position, lastTradedPrice float64) float64 {
	if position.Direction == LONG {
		return float64(position.Size) * (lastTradedPrice - position.EntryPrice)
	}
	return float64(position.Size) * (position.EntryPrice - lastTradedPrice)
}

//USDMarginedLiquidationPrice calculates the liquidation price for the positions when trading USD margined assets.
//This calculation assumes the isolated trading position mode
//More info on https://help.bybit.com/hc/en-us/articles/900000181046-Liquidation-Price-USDT-Contract
func USDMarginedLiquidationPrice(position *Position) float64 {
	leverage := math.Max(1.0, float64(position.Leverage))
	IMR := 1.0 / leverage
	if position.Direction == LONG {
		return position.EntryPrice * (1 - IMR + MMR)
	}
	return position.EntryPrice * (1 + IMR - MMR)
}

//CoinMarginedLiquidationPrice calculates the liquidation price for the position when trading COIN margined assets
//This calculation assumes the isolated trading position mode
//More info on https://help.bybit.com/hc/en-us/articles/360039261334-How-to-calculate-Liquidation-Price-Inverse-Contract
func CoinMarginedLiquidationPrice(position *Position) float64 {
	leverage := math.Max(1.0, float64(position.Leverage))
	if position.Direction == LONG {
		return (position.EntryPrice * leverage) / (leverage + 1 - (MMR * leverage))
	}
	return (position.EntryPrice * leverage) / (leverage - 1 + (MMR * leverage))
}
