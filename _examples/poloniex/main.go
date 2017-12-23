package main

import (
	"log"
  "time"
	poloClient "github.com/dfontana/Cryptical/poloniex"
	"github.com/dfontana/Cryptical/plot"
)

func main() {
		// Get historic data and save to CSV
		recsP, err := poloClient.Historic("USDT_ETH", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 300)
		if err != nil {
			log.Fatal(err)
		}
		poloClient.CSV("./outP.csv", recsP)
		
		// Create and save a bollinger plot
		bollinger()
	
		// Create and save a MACD plot
		macd()
	
		// Stream live data for 10 seconds
		stream()
}

func bollinger() {
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

func macd(){
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

func stream(){
		// Asynchronously fetch data to messages channel.
		messages := make(chan poloClient.WSOrderbook)
		quit := make(chan bool)
		go poloClient.Live([]string{"USDT_BTC", "USDT_ETH"}, messages, quit)
	
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