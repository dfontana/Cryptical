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
 * Example where we build a Bollinger Band plot from historic data
 */
func gdaxBollinger() {
	daysBack := 150

	var records []gdax.HistoricRate
	start := time.Now().AddDate(0, 0, -daysBack)
	end := time.Now()
	gran := 24 * 60 * 60
	expected := int(math.Ceil(end.Sub(start).Seconds()/float64(gran))) + 1

	// Repetition due to instability in GDAX api
	for {
		records = gdaxClient.Historic("ETH-USD", start, end, gran)
		log.Printf("Data returned from API: %d/%d\n", len(records), expected)
		if len(records) == expected {
			break // Correct amount of data found
		}
		time.Sleep(time.Duration(3) * time.Second)
	}

	// Reduce to array of values & their times
	hist := make([]plot.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = plot.TimeSeries{
			val.Time,
			val.High,
		}
	}

	// Make a Bollinger Plot
	s := time.Now()
	comp := plot.Bollinger{}
	comp.Compute(hist, 20)
	comp.Plot("./testbb.png")
	e3 := time.Since(s)

	log.Printf("Bollinger Plot Took: %s", e3)
}

/**
* Example where we build a Bollinger Band plot from historic data
*/
func polBollinger() {
daysBack := 150

// Past 150 Daily records
recsP, err := poloClient.Historic("USDT_ETH", time.Now().AddDate(0, 0, -daysBack), time.Now(), 1800)
if err != nil {
	log.Fatal(err)
}

// Reduce to array of values & their times
hist := make([]plot.TimeSeries, len(recsP))
for i, val := range recsP {
	hist[i] = plot.TimeSeries{
		val.Date.Time,
		val.High,
	}
}

// Make a Bollinger Plot
s := time.Now()
comp := plot.Bollinger{}
comp.Compute(hist, 20)
comp.Plot("./testbb.png")
e3 := time.Since(s)

log.Printf("Bollinger Plot Took: %s", e3)
}