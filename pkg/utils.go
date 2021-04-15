package pkg

import "math"

//MMR - maintenance margin rate
const MMR = 0.005

func stdDev(numbers []float64) float64 {
	total := 0.0
	mean := mean(numbers)
	for _, number := range numbers {
		total += math.Pow(number-mean, 2)
	}
	variance := total / float64(len(numbers)-1)
	return math.Sqrt(variance)
}

func mean(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers)-1)
}

func sharpe(investmentReturn, riskFreeReturn, stdDev float64) float64 {
	return (investmentReturn - riskFreeReturn) / stdDev
}
