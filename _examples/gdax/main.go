package main

import (
	"log"
	"math"
	"time"

	"github.com/dfontana/Cryptical/computation"
	gdaxClient "github.com/dfontana/Cryptical/gdax"
	gdax "github.com/preichenberger/go-gdax"
)

func main() {
	// Get historic data and save to CSV
	recs := gdaxClient.Historic("ETH-USD", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 200)
	gdaxClient.CSV("./outG.csv", recs)

	// Create and save a bollinger plot
	bollinger()

	// Create and save a MACD plot
	macd()

	// Stream live data for 10 seconds
	stream()
}

func bollinger() {
	daysBack := 20

	var records []gdax.HistoricRate
	start := time.Now().AddDate(0, 0, -daysBack)
	end := time.Now()
	gran := 1800
	expected := 959 // Hard code

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
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Time,
			val.High,
		}
	}

	// Make a Bollinger Plot
	s := time.Now()
	comp := computation.Bollinger{}
	comp.Compute(hist, 20)
	comp.Plot("./testbb.png")
	e3 := time.Since(s)

	log.Printf("Bollinger Plot Took: %s", e3)
}

func macd() {
	daysBack := 150

	// Past 150 days for ETH daily.
	s := time.Now()
	var records []gdax.HistoricRate
	start := time.Now().AddDate(0, 0, -daysBack)
	end := time.Now()
	gran := 24 * 60 * 60
	expected := int(math.Ceil(end.Sub(start).Seconds() / float64(gran)))
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
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Time,
			val.High,
		}
	}
	e2 := time.Since(s)

	// MACD: 12 fast, 26 slow, 9 signal
	s = time.Now()
	comp := computation.MACD{}
	if err := comp.Compute(hist, 12, 26, 9); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./testmacd.png")
	e3 := time.Since(s)

	log.Printf("Timings:\n\tHistory: %s\n\tTimeSeries: %s\n\tMACD: %s", e1, e2, e3)
}

func stream() {
	// Asynchronously fetch data to messages channel.
	messages := make(chan gdaxClient.WsMatch)
	quit := make(chan bool)
	go gdaxClient.Live([]string{"ETH-USD", "BTC-USD"}, messages, quit)

	// Kill the livefeed after 10 seconds.
	go func() {
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	// Loop until something stops the socket feed (error or disabled)
	for msg := range messages {
		log.Printf("%+v\n", msg)
	}
}
