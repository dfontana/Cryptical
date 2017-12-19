package common

import (
	"errors"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"os"
	"image/png"
	"time"
	"log"
)

// MACD Describes a single MACD computation. Should call Populate() to fill
// this structs fields.
type MACD struct {
	Entries	[]MACDItem	// Time series data for MACD
}

// MACDIten describes a single entry of MACD; its values for a given
// date.
type MACDItem struct {
	Time		int64		// Epoch Time
	MACD		float64	// MACD value
	Signal	float64	// Signal value
}

type TimeSeries struct {
	Time	int64		// Epoch time of datapoint
	Data 	float64	// Data point 
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
	m.Entries = make([]MACDItem, len(sign))
	for i,_ := range sign {
		m.Entries[i] = MACDItem {
			sign[i].Time,
			macd[i].Data,
			sign[i].Data,
		}
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

func (m *MACD) Plot() error{
	if m.Entries == nil {
		return errors.New("Nothing to plot, did you Populate() your data?")
	}
	xv := make([]time.Time, len(m.Entries))
	yMv := make([]float64, len(m.Entries))
	ySv := make([]float64, len(m.Entries))

	for i,entry := range m.Entries {
		xv[i] = time.Unix(entry.Time, 0)
		yMv[i] = entry.MACD
		ySv[i] = entry.Signal
	}


	mSeries := chart.TimeSeries{
		Name: "MACD",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorBlue,
		},
		XValues: xv,
		YValues: yMv,
	}

	sSeries := chart.TimeSeries{
		Name: "Signal",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorRed,
		},
		XValues: xv,
		YValues: ySv,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Style: chart.Style{Show: true},
			Range: &chart.ContinuousRange{
				Max: 80.0,
				Min: -30.0,
			},
		},
		Series: []chart.Series{
			mSeries,
			sSeries,
		},
	}

	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)
	image, err := collector.Image()
	if err != nil {
		log.Fatal(err)
	}
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		// Handle error
		log.Fatal(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, image)

	// Don't forget to close files
	outputFile.Close()

	return nil
}
