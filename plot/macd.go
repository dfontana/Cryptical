package common

import (
	"errors"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"image/png"
	"time"
)

// MACD Describes a single MACD computation. Should call Populate() to fill
// this structs fields.
type MACD struct {
	Time	 	[]time.Time // Array of times
	MACD		[]float64		// Array of MACD values, corresponding to time.
	Signal	[]float64 	// Array of Signal Values, corresponding to time.
	Hist		[]float64
}

// 	Populate will fill in  computes Moving Average Convergence Divergence.
//	closingPrices:
//		An array of closing prices associated with times. This doesn't have to
//		be daily.
// 	fast, slow, signal: 
//		Values to operate on closingPrices - assumed to be in the same time
//		unit closingPrices was passed as. For days this is typically 12,26,9.
// 	Returns: 
//		Error is one occured, otherwise the Entries of this struct are now filled.
func (m *MACD) Populate(closingPrices []TimeSeries, fast, slow, signal int) (error) {
	if fast >= slow {
		return errors.New("Fast > slow. No.")
	}

	// Calculate EMAs(t)
	emaFast, err := ema(closingPrices, fast)
	if err != nil {
		return err
	}

	emaSlow, err := ema(closingPrices, slow)
	if err != nil {
		return err
	}

	// Calculate MACD(t)
	for i, _ := range emaSlow {
		emaSlow[i] = TimeSeries{
			emaFast[i+(slow-fast)].Time,
			emaFast[i+(slow-fast)].Data - emaSlow[i].Data,
		}
	}
	macd := emaSlow
	macd = macd[signal:] // Trim burned data from signal calc

	// Calculate signal
	sign, err := ema(macd, signal)
	if err != nil {
		return err
	}

	// Join our data into MACD items, then into an MACD struct
	m.Time = make([]time.Time, len(sign))
	m.MACD = make([]float64, len(sign))
	m.Signal = make([]float64, len(sign))
	m.Hist = make([]float64, len(sign))
	for i,_ := range sign {
		m.Time[i] = sign[i].Time
		m.MACD[i] = macd[i].Data
		m.Signal[i] = sign[i].Data
		m.Hist[i] = macd[i].Data - sign[i].Data
	}
	return nil
}

// EMA computes Exponential Moving Average for given period within the given
// slice. Returns array of values - ema per time period.
func ema(closingPrices []TimeSeries, period int) ([]TimeSeries, error) {
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
	result := make([]TimeSeries, len(validHist))
	for k, price := range validHist {
		prevEMA = multi*price.Data + (1-multi)*prevEMA
		result[k] = TimeSeries {price.Time, prevEMA}
	}

	// The expected result is of len closingPrices - period.
	return result, nil
}

// sma is the Simple Moving Average for given slice
func sma(closingPrices []TimeSeries) float64 {
	if len(closingPrices) == 0 {
		return 0
	}

	var sum float64
	for _, v := range closingPrices {
		sum += v.Data
	}

	return sum / float64(len(closingPrices))
}

// Plot creates a plot from this computation instance, saved at the 
// given path.
func (m *MACD) Plot(path string) error{
	if m.Time == nil || m.MACD == nil || m.Signal == nil {
		return errors.New("Nothing to plot, did you Populate() your data?")
	}

	mSeries := chart.TimeSeries{
		Name: "MACD",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorBlue,
			DotColor: drawing.ColorBlue,
			DotWidth:	2.0,
		},
		XValues: m.Time,
		YValues: m.MACD,
	}

	sSeries := chart.TimeSeries{
		Name: "Signal",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorRed,
			DotColor: drawing.ColorRed,
			DotWidth:	2.0,
		},
		XValues: m.Time,
		YValues: m.Signal,
	}

	hSeries := chart.TimeSeries{
		Name: "HistDiff",
		XValues: m.Time,
		YValues: m.Hist,
	}

	histSeries := chart.HistogramSeries{
		Name: "Hist",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("426993"),
			FillColor: drawing.ColorFromHex("426993").WithAlpha(64),
		},
		YAxis: chart.YAxisPrimary,
		InnerSeries: hSeries,
	}

	// Figure out our Y Bounds real quick:
	mMin, mMax := minMax(m.MACD)
	sMin, sMax := minMax(m.Signal)
	lower := sMin - 50
	if mMin < sMin {
		lower = mMin - 50
	}
	upper := sMax + 50
	if mMax > sMax {
		upper = mMax + 50
	}
	
	// Now create the chart
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
			mSeries,
			sSeries,
			histSeries,
		},
	}

	if err := saveImage(graph, path); err != nil {
		return err
	}

	return nil
}

// Renders and saves graph.
func saveImage(graph chart.Chart, path string) error {
	// Write image to buffer
	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)
	image, err := collector.Image()
	if err != nil {
		return err
	}
	
	// Save buffer to file (after encoding)
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	png.Encode(outputFile, image)
	outputFile.Close()
	return nil
}

// Returns the min and max from a slice
func minMax(vals []float64) (float64, float64) {
	min := vals[0]
	max := vals[0]
	for _, val := range(vals){
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return min, max
}
