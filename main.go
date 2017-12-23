package main

import (
  "log"
  "math"
  "time"

  gdax "github.com/preichenberger/go-gdax"
  
  "github.com/dfontana/Cryptical/plot"
  gdaxClient "github.com/dfontana/Cryptical/gdax"
  poloClient "github.com/dfontana/Cryptical/poloniex"
)

// The main function will always execute when called from "go run". So all we do
// here is call the example we'd like to see
func main() {
  log.Println(BreakEven("ETH-USD", 25.0))
}


func inference(startDate time.Time, soldprices chan float64) {

	// These are in 5 min periods meaining we look back 12, 5, and 3.5 hours
	slow := 24 * 6
	fast := 10 * 6
	sign := 7 * 6
	gran := 300 // Seconds

	// Will probably want ~3 * Slow data points to ensure we have enough.
	startHist := startDate.Add(-3 * slow * gran * time.Second) 
  records, err := poloClient.Historic("USDT_ETH", startHist, startDate, gran)
  if err != nil {
    log.Fatal(err)
  }

  // Convert said data to the needed format to run our MACD computation.
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
