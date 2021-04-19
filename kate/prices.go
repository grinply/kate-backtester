package kate

//DataPoint is a unit that encapsulates OHLCV price data
type DataPoint struct {
	Event
	open, high, low, close, volume float64
}

//OHLCV - Represents a datapoint in candle format that contain Open,High, Low, Close prices and Volume data
type OHLCV interface {
	Close() float64
	Open() float64
	High() float64
	Low() float64
	Volume() float64
}

func (candle DataPoint) Open() float64 {
	return candle.open
}

func (candle DataPoint) High() float64 {
	return candle.high
}

func (candle DataPoint) Low() float64 {
	return candle.low
}

func (candle DataPoint) Close() float64 {
	return candle.close
}

func (candle DataPoint) Volume() float64 {
	return candle.volume
}
