package examples

import (
  "log"
  "math"
  "time"

  gdax "github.com/preichenberger/go-gdax"
  
  "github.com/dfontana/Cryptical/plot"
  gdaxClient "github.com/dfontana/Cryptical/gdax"
  poloClient "github.com/dfontana/Cryptical/poloniex"
)

/**
 * In this example, we grab historic data, transform it to be compatible with
 * MACD, compute the MACD of this data and then plot it.
 */
func gdaxMACD() {
  daysBack := 150

  // Past 150 days for ETH daily.
  s := time.Now()
  var records []gdax.HistoricRate
  start := time.Now().AddDate(0, 0, -daysBack)
  end := time.Now()
  gran := 24 * 60 * 60
  expected := int(math.Ceil(end.Sub(start).Seconds()/float64(gran))) + 1
  for {
    records = gdaxClient.Historic("ETH-USD", start, end, gran)
    log.Printf("Data returned from API: %d/%d\n", len(records), expected)
    if len(records) == expected {
      break // Correct amount of data found
    }
    time.Sleep(time.Duration(3) * time.Second)
  }
  e1 := time.Since(s)

  s = time.Now()
  // Reduce to array of close values & their times
  hist := make([]plot.TimeSeries, len(records))
  for i, val := range records {
    hist[i] = plot.TimeSeries{
      val.Time,
      val.High,
    }
  }
  e2 := time.Since(s)

  // MACD: 12 fast, 26 slow, 9 signal
  s = time.Now()
  comp := plot.MACD{}
  if err := comp.Compute(hist, 12, 26, 9); err != nil {
    log.Fatal(err)
  }
  comp.Plot("./testmacd.png")
  e3 := time.Since(s)

  log.Printf("Timings:\n\tHistory: %s\n\tTimeSeries: %s\n\tMACD: %s", e1, e2, e3)
}

func polMACD() {
  daysBack := 150

  // Past 150 days for ETH daily.
  s := time.Now()
  start := time.Now().AddDate(0, 0, -daysBack)
  end := time.Now()
  gran := 24 * 60 * 60
  records, err := poloClient.Historic("USDT_ETH", start, end, gran)
  if err != nil {
    log.Fatal(err)
  }
  e1 := time.Since(s)

  s = time.Now()
  // Reduce to array of close values & their times
  hist := make([]plot.TimeSeries, len(records))
  for i, val := range records {
    hist[i] = plot.TimeSeries{
      val.Date.Time,
      val.High,
    }
  }
  e2 := time.Since(s)

  // MACD: 12 fast, 26 slow, 9 signal
  s = time.Now()
  comp := plot.MACD{}
  if err := comp.Compute(hist, 12, 26, 9); err != nil {
    log.Fatal(err)
  }
  comp.Plot("./test.png")
  e3 := time.Since(s)

  log.Printf("Timings:\n\tHistory: %s\n\tTimeSeries: %s\n\tMACD: %s", e1, e2, e3)
}