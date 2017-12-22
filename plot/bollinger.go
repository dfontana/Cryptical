package plot

import (
  "errors"
  "github.com/wcharczuk/go-chart"
  "time"
  "math"
)

// Bollinger handles generating bollinger band plots from the given historical
// price data. This data should be in the form of time series structs, to
// handle time senstive data. Populate the struct, then call its methods.
type Bollinger struct {
  History	[]TimeSeries	// Array of historical data points used for computation
  time		[]time.Time		// times for each data point in upper,sma,lower
  upper		[]float64			// Simple Moving Average + 2*STDDev
  sma     []float64			// Simple Moving Average
  lower		[]float64			// Simple Moving Average - 2*STDDev
}

// Compute determines the bollinger bands of this structs history
// This function expects the history field to have been filled, and will error
// if omitted.
// avgLen: The number of data points to compute the Simple Moving Average over.
func (b *Bollinger) Compute(hist []TimeSeries, avgLen int) error {
  if len(hist) < avgLen-1 {
    return errors.New("Not enough data to compute the requested SMA timeframe.")
  }

  b.History = hist

  // Compute average and stdDev, filling in the bands
  b.time = make([]time.Time, len(b.History)-avgLen+1)
  b.upper = make([]float64, len(b.History)-avgLen+1)
  b.sma = make([]float64, len(b.History)-avgLen+1)
  b.lower = make([]float64, len(b.History)-avgLen+1)
  for i,_ := range b.sma {
    sigma, mu := sigmaMu(b.History[i:i+avgLen])
    b.time[i] = b.History[i+avgLen-1].Time 
    b.upper[i] = mu + 2*sigma
    b.sma[i] = mu
    b.lower[i] = mu - 2*sigma
  }
  return nil
}

// Plot will create a bollinger plot from data stored in the type,
// saved to the given path. Since Go-Chart provides a means to compute
// this itself, you don't have to call populate on this.
func (b *Bollinger) Plot(path string) error {
  if b.time == nil || b.History == nil{
    return errors.New("No data to plot, please compute first.")
  }

  // Reduce the time series into just x's and y's
  drawable := len(b.time)
  xv := make([]time.Time, drawable)
  yv := make([]float64, drawable)
  for i,item := range b.History[len(b.History)-drawable:] {
    xv[i] = item.Time
    yv[i] = item.Data
  }

  // Construct series for plotting.
  hSeries := chart.TimeSeries {
    Name: "Prices",
    Style: chart.Style{
      Show: true,
      StrokeColor: chart.ColorAlternateBlue,
      FillColor: chart.ColorAlternateBlue.WithAlpha(70),
    },
    XValues: xv,
    YValues: yv,
  }

  uSeries := chart.TimeSeries {
    Name: "Prices",
    Style: chart.Style{
      Show: true,
      StrokeColor: chart.ColorAlternateGreen,
    },
    XValues: xv,
    YValues: b.upper,
  }

  sSeries := chart.TimeSeries {
    Name: "Prices",
    Style: chart.Style{
      Show: true,
      StrokeColor: chart.ColorAlternateGreen,
      StrokeDashArray: []float64{0.5, 0.5},
    },
    XValues: xv,
    YValues: b.sma,
  }

  lSeries := chart.TimeSeries {
    Name: "Prices",
    Style: chart.Style{
      Show: true,
      StrokeColor: chart.ColorAlternateGreen,
    },
    XValues: xv,
    YValues: b.lower,
  }

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
      hSeries,
      sSeries,
      uSeries,
      lSeries,
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
  sum := 0.0
  for _,x := range(data) {
    sum += math.Pow(x.Data - mu, 2)
  }
  sigma := math.Sqrt((1 / float64(n)) * sum)

  return sigma, mu
}