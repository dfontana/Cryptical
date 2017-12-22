package plot

import (
	"errors"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"time"
	"math"
)

// Bollinger handles generating bollinger band plots from the given historical
// price data. This data should be in the form of time series structs, to
// handle time senstive data. Populate the struct, then call its methods.
type Bollinger struct {
	History	[]TimeSeries	// Array of historical data points used for computation
	upper		[]TimeSeries	// Simple Moving Average + 2*STDDev
	sma     []TimeSeries	// Simple Moving Average
	lower		[]TimeSeries	// Simple Moving Average - 2*STDDev
}

// Compute determines the bollinger bands of this structs history
// This function expects the history field to have been filled, and will error
// if omitted.
// avgLen: The number of data points to compute the Simple Moving Average over.
func (b *Bollinger) Compute(avgLen int) error {
	if b.History == nil {
		return errors.New("No data to plot, please fill History field.")
	}

	if len(b.History) < avgLen-1 {
		return errors.New("Not enough data to compute the requested SMA timeframe.")
	}

	// Compute average and stdDev, filling in the bands
	b.upper = make([]TimeSeries, len(b.History)-avgLen+1)
	b.sma = make([]TimeSeries, len(b.History)-avgLen+1)
	b.lower = make([]TimeSeries, len(b.History)-avgLen+1)
	for i,_ := range sma {
		sigma, mu := sigmaMu(b.History[i:i+avgLen])
		b.upper[i] = mu + 2*sigma
		b.sma[i] = mu
		b.lower[i] = mu - 2*sigma
	}
}

// Plot will create a bollinger plot from data stored in the type,
// saved to the given path. Since Go-Chart provides a means to compute
// this itself, you don't have to call populate on this.
func (b *Bollinger) Plot(path string) error {
	if b.History == nil {
		return errors.New("No data to plot, please fill History field.")
	}

	// Reduce the time series into just x's and y's
	xv := make([]time.Time, len(b.History))
	yv := make([]float64, len(b.History))
	for i,item := range b.History {
		xv[i] = item.Time
		yv[i] = item.Data
	}

	// Construct series for plotting.
	hSeries := chart.TimeSeries {
		Name: "Prices",
		Style: chart.Style{
			Show: true,
			StrokeColor: chart.ColorBlue,
		},
		XValues: xv,
		YValues: yv,
	}

	// Build Bollinger Bands
	bbSeries := &chart.BollingerBandsSeries {
		Name: "Bollinger",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("426993"),
			FillColor:   drawing.ColorFromHex("426993").WithAlpha(64),
		},
		InnerSeries: hSeries,
	}

	// Figure out our Y Bounds real quick:
	hMin, hMax := MinMax(yv)
	lower := hMin - 50
	upper := hMax + 50

	// Plot it!
	graph := chart.Chart{
		Width: 1920,
		Height: 1080,
		DPI: 100,
		XAxis: chart.XAxis{
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
			ValueFormatter: chart.TimeValueFormatter,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
			Range: &chart.ContinuousRange{
				Max: upper,
				Min: lower,
			},
		},
		Series: []chart.Series{
			bbSeries,
			hSeries,
		},
	}

	if err := SaveImage(graph, path); err != nil {
		return err
	}

	return nil
}

// sigmaMu returns the standard deviation (sigma) and average
// (mu) of the slice provided to it.
func sigmaMu(data []TimeSeries) (float64, float64) {
	n := len(data)
	if n == 0 {
		return 0, 0
	}

	mu := Sma(data)
	sum := 0
	for _,x := range(data) {
		sum += math.Pow(x - mu, 2)
	}
	sigma := math.Sqrt((1 / n) * sum)

	return sigma, mu
}