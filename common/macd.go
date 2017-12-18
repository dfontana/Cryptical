package common

import (
	"errors"
)

// MACD computes Moving Average Convergence Divergence.
// hist: historic prices for a length of time. This can be days, hours, etc
//		 - each index in this array represents a value for a time period.
//		   They should be in increasing order
// fast, slow, signal: Values to operate on hist - assumed to be in the same
//		 time unit hist was passed as. For days this is typically 12,26,9.
// Returns MACD(t), Signal(t), and error.
func MACD(hist []float64, fast, slow, signal int) ([]float64, []float64, error) {
	if fast >= slow {
		return nil, nil, errors.New("Fast > slow. No.")
	}

	// Calculate EMAs(t)
	emaFast, err := ema(hist, fast)
	if err != nil {
		return nil, nil, err
	}

	emaSlow, err := ema(hist, slow)
	if err != nil {
		return nil, nil, err
	}

	// Calculate MACD(t)
	for i, _ := range emaSlow {
		emaSlow[i] = emaFast[i+(slow-fast)] - emaSlow[i]
	}
	macd := emaSlow

	// Calculate signal
	sign, err := ema(macd, signal)
	if err != nil {
		return nil, nil, err
	}

	return macd, sign, nil
}

// EMA computes Exponential Moving Average for given period within the given
// slice. Returns array of values - ema per time period.
func ema(hist []float64, period int) ([]float64, error) {
	if len(hist) < period {
		return nil, errors.New("Need more history.")
	}
	// Starting point is a simple average
	prevEMA := sma(hist[0:period])

	// From starting point onward we start getting valid EMAs
	emaHist := hist[period:len(hist)]
	multi		:= 2 / float64(period+1)
	result := make([]float64, len(emaHist)+1)
	result[0] = prevEMA

	for k, price := range emaHist {
		prevEMA = multi*price + (1-multi)*prevEMA
		result[k+1] = prevEMA
	}

	return result, nil
}

// sma is the Simple Moving Average for given slice
func sma(hist []float64) float64 {
	if len(hist) == 0 {
		return 0
	}

	var sum float64
	for _, v := range hist {
		sum += v
	}

	return sum / float64(len(hist))
}
