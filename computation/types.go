package computation

import "time"

// TimeSeries describes a single data point that should be paired with a point
// in time.
type TimeSeries struct {
	Time time.Time // Time of data point
	Data float64   // Data point
}

// TradeString defines constants Buy, Sell, and Hodl.
type TradeString string

const (
	// Buy - Should submit a buy action
	Buy TradeString = "buy"

	// Sell - Should submit a sell action
	Sell TradeString = "sell"

	// Hodl - Should submit a hodl action
	Hodl TradeString = "hodl"
)

// Trade resembles a simulation trade, of bought/sold and the amount in ETH/USD
type Trade struct {
	Type   TradeString
	Crypto float64
	USD    float64
}

// Computation defines a contract for types that is able to be used in inference.
type Computation interface {
	AddPoint(t TimeSeries) error
	Analyze() (Trade, error)
	Compute() error
	Plot(path string) error
}
