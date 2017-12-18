package common

import (
	"errors"
)

// MACD computes Moving Average Convergence Divergence.
// closingPrices: historic prices for a length of time. This can be days, hours, etc
//		 - each index in this array represents a value for a time period.
//		   They should be in increasing order
// fast, slow, signal: Values to operate on closingPrices - assumed to be in the same
//		 time unit closingPrices was passed as. For days this is typically 12,26,9.
// Returns MACD(t), Signal(t), and error.
func MACD(closingPrices []float64, fast, slow, signal int) ([]float64, []float64, error) {
	if fast >= slow {
		return nil, nil, errors.New("Fast > slow. No.")
	}

	// Calculate EMAs(t)
	emaFast, err := ema(closingPrices, fast)
	if err != nil {
		return nil, nil, err
	}

	emaSlow, err := ema(closingPrices, slow)
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
func ema(closingPrices []float64, period int) ([]float64, error) {
	if len(closingPrices) < period {
		return nil, errors.New("Need more history.")
	}
	// Starting point is a simple average
	prevEMA := sma(closingPrices[0:period])

	// Truncate the first period of days off the history, since those are
	// are used to initialize the prevEMA
	validHist := closingPrices[period:len(closingPrices)]
	multi		:= 2 / float64(period+1)

	// Only store the valid EMAs -> Compute them.
	result := make([]float64, len(validHist))
	for k, price := range validHist {
		prevEMA = multi*price + (1-multi)*prevEMA
		result[k] = prevEMA
	}

	// The expectation is result is of len closingPrices - period.
	return result, nil
}

// sma is the Simple Moving Average for given slice
func sma(closingPrices []float64) float64 {
	if len(closingPrices) == 0 {
		return 0
	}

	var sum float64
	for _, v := range closingPrices {
		sum += v
	}

	return sum / float64(len(closingPrices))
}
