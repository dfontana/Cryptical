package computation

import (
	"errors"
	"time"

	chart "github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// MACD Describes a single MACD computation. Should call Populate() to fill
// this structs fields.
type MACD struct {
	// Model Paramters
	Data []TimeSeries
	Fast int
	Slow int
	Sign int

	// Computed Model
	Time   []time.Time // Array of times
	MACD   []float64   // Array of MACD values, corresponding to time.
	Signal []float64   // Array of Signal Values, corresponding to time.
	Hist   []float64   // Historgram values

	// Trading Strategy based on model
	Strategy func(m *MACD) Trade
}

// AddPoint inserts a new data point to the already computed model. If the
// model was not computed then an error is returned.
func (m *MACD) AddPoint(t TimeSeries) error {
	if len(m.MACD) == 0 {
		return errors.New("initial model must be made before adding new data points")
	}
	m.Data = append(m.Data, t)
	return nil
}

// Analyze exectues this model's strategy. If the strategy is not defined, then
// an error is thrown
func (m *MACD) Analyze() (Trade, error) {
	if m.Strategy == nil {
		return Trade{}, errors.New("model does not have a strategy defined")
	}
	return m.Strategy(m), nil
}

// Compute populates the defined model based on the given parameters. If any
// values are missing (or a problem occurs) then an error is returned.
func (m *MACD) Compute() error {
	if len(m.Data) == 0 || m.Fast == 0 || m.Slow == 0 || m.Sign == 0 {
		return errors.New("missing parameters in model definition")
	}

	if m.Fast >= m.Slow {
		return errors.New("fast > slow in model, I can't work like this")
	}

	// Calculate EMAs(t)
	emaFast, err := ema(m.Data, m.Fast)
	if err != nil {
		return err
	}

	emaSlow, err := ema(m.Data, m.Slow)
	if err != nil {
		return err
	}

	// Calculate MACD(t)
	for i := range emaSlow {
		emaSlow[i] = TimeSeries{
			emaFast[i+(m.Slow-m.Fast)].Time,
			emaFast[i+(m.Slow-m.Fast)].Data - emaSlow[i].Data,
		}
	}
	macd := emaSlow
	macd = macd[m.Sign:] // Trim burned data from signal calc

	// Calculate signal
	sign, err := ema(macd, m.Sign)
	if err != nil {
		return err
	}

	// Join our data into MACD items, then into an MACD struct
	m.Time = make([]time.Time, len(sign))
	m.MACD = make([]float64, len(sign))
	m.Signal = make([]float64, len(sign))
	m.Hist = make([]float64, len(sign))
	for i := range sign {
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
		return nil, errors.New("need more history")
	}
	// Starting point is a simple average
	prevEMA := Sma(closingPrices[0:period])

	// Truncate the first period of days off the history, since those are
	// are used to initialize the prevEMA
	validHist := closingPrices[period:len(closingPrices)]
	multi := 2 / float64(period+1)

	// Only store the valid EMAs -> Compute them.
	result := make([]TimeSeries, len(validHist))
	for k, price := range validHist {
		prevEMA = multi*price.Data + (1-multi)*prevEMA
		result[k] = TimeSeries{price.Time, prevEMA}
	}

	// The expected result is of len closingPrices - period.
	return result, nil
}

// Plot creates a plot from this computation instance, saved at the
// given path.
func (m *MACD) Plot(path string) error {
	if m.Time == nil || m.MACD == nil || m.Signal == nil {
		return errors.New("Nothing to plot, did you Populate() your data?")
	}

	mSeries := chart.TimeSeries{
		Name: "MACD",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorBlue,
			DotColor:    drawing.ColorBlue,
			DotWidth:    2.0,
		},
		XValues: m.Time,
		YValues: m.MACD,
	}

	sSeries := chart.TimeSeries{
		Name: "Signal",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorRed,
			DotColor:    drawing.ColorRed,
			DotWidth:    2.0,
		},
		XValues: m.Time,
		YValues: m.Signal,
	}

	hSeries := chart.TimeSeries{
		Name:    "HistDiff",
		XValues: m.Time,
		YValues: m.Hist,
	}

	histSeries := chart.HistogramSeries{
		Name: "Hist",
		Style: chart.Style{
			Show:        true,
			StrokeColor: drawing.ColorFromHex("426993"),
			FillColor:   drawing.ColorFromHex("426993").WithAlpha(64),
		},
		YAxis:       chart.YAxisPrimary,
		InnerSeries: hSeries,
	}

	// Figure out our Y Bounds real quick:
	mMin, mMax := MinMax(m.MACD)
	sMin, sMax := MinMax(m.Signal)
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
		Width:  1920,
		Height: 1080,
		DPI:    100,
		XAxis: chart.XAxis{
			Style:          chart.Style{Show: true},
			TickPosition:   chart.TickPositionBetweenTicks,
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

	err := SaveImage(graph, path)
	return err
}
