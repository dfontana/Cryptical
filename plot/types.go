package plot

import "time"

// TimeSeries describes a single data point that should be paired with a point
// in time.
type TimeSeries struct {
	Time	time.Time		// Time of data point
	Data 	float64			// Data point 
}