package computation

import "time"

// TimeSeries describes a single data point that should be paired with a point
// in time.
type TimeSeries struct {
	Time time.Time // Time of data point
	Data float64   // Data point
}

// TradeString defines constants Buy, Sell, and Hodl.
type ActionString string

const (
	// Buy - Should submit a buy action
	Buy ActionString = "buy"

	// Sell - Should submit a sell action
	Sell ActionString = "sell"

	// Hodl - Should submit a hodl action
	Hodl ActionString = "hodl"
)

// Action describes what to do, how much cryptocurrency to move, and the current value of 1 of those
// Cryptocurrencies in this bot's pairing
type Action struct {
	Type   ActionString
	Crypto float64
	Value  float64
}

// Computation defines a contract for types that is able to be used in inference.
type Computation interface {
	AddPoint(t TimeSeries) error
	Analyze() (Action, error)
	Compute() error
	Plot(path string) error
}
