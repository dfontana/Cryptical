package main

import (
	"log"
	"math"
	"time"

	gdax "github.com/preichenberger/go-gdax"
	
	"./common"
	gdaxClient "./gdax"
	poloClient "./poloniex"
)

func main() {
	polMACD()
}

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
	hist := make([]common.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = common.TimeSeries{
			val.Time,
			val.High,
		}
	}
	e2 := time.Since(s)

	// MACD: 12 fast, 26 slow, 9 signal
	s = time.Now()
	comp := common.MACD{}
	if err := comp.Populate(hist, 12, 26, 9); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./test.png")
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
	hist := make([]common.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = common.TimeSeries{
			val.Date.Time,
			val.High,
		}
	}
	e2 := time.Since(s)

	// MACD: 12 fast, 26 slow, 9 signal
	s = time.Now()
	comp := common.MACD{}
	if err := comp.Populate(hist, 12, 26, 9); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./test.png")
	e3 := time.Since(s)

	log.Printf("Timings:\n\tHistory: %s\n\tTimeSeries: %s\n\tMACD: %s", e1, e2, e3)
}

func polHist() {
	recsP, err := poloClient.Historic("USDT_ETH", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 300)
	if err != nil {
		log.Fatal(err)
	}
	poloClient.CSV("./outP.csv", recsP)
}

func gdaxHist() {
	recs := gdaxClient.Historic("ETH-USD", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 200)
	gdaxClient.CSV("./outG.csv", recs)
}

func polLive() {
		// Asynchronously fetch data to messages channel.
		messages := make(chan poloClient.WSOrderbook)
		quit := make(chan bool)
		go poloClient.Live(messages, quit)
	
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

func gdaxLive() {
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
