package main

import (
	"log"
	"time"

	"./common"
	gdax "./gdax"
	pol "./poloniex"
)

// The main function will always execute when called from "go run". So all we do
// here is call the example we'd like to see
func main() {
	gdaxBollinger()
}


/**
 * Example where we build a Bollinger Band plot from historic data
 */
func gdaxBollinger() {
		g := gdax.GDAX{[]string{"ETH-USD"}}
	
		daysBack := 150

		// Past 150 days for ETH daily.
		records := g.Historic(g.Currencies[0], time.Now().AddDate(0, 0, -daysBack), time.Now(), 24*60*60)

		// Due to unreliability in gdax API, we have to check if more data was returned than requested.
		if len(records) > daysBack+1 {
			log.Fatalf("GDAX API gave too many records: %d/%d", len(records), daysBack+1)
		}

		// Reduce to array of close values & their times
		hist := make([]common.TimeSeries, len(records))
		for i, val := range records {
			hist[i] = common.TimeSeries{
				time.Unix(val.Time, 0),
				val.High,
			}
		}

		// Make a Bollinger Plot
		s := time.Now()
		comp := common.Bollinger{hist}
		comp.Plot("./testbb.png")
		e3 := time.Since(s)

		log.Printf("Bollinger Plot Took: %s", e3)
}

/**
 * In this example, we grab historic data, transform it to be compatible with
 * MACD, compute the MACD of this data and then plot it.
 */
func gdaxMACD() {
	g := gdax.GDAX{[]string{"ETH-USD"}}

	daysBack := 150

	// Past 150 days for ETH daily.
	s := time.Now()
	records := g.Historic(g.Currencies[0], time.Now().AddDate(0, 0, -daysBack), time.Now(), 24*60*60)
	e1 := time.Since(s)

	// Due to unreliability in gdax API, we have to check if more data was returned than requested.
	if len(records) > daysBack+1 {
		log.Fatalf("GDAX API gave too many records: %d/%d", len(records), daysBack+1)
	}

	s = time.Now()
	// Reduce to array of close values & their times
	hist := make([]common.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = common.TimeSeries{
			time.Unix(val.Time, 0),
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
	comp.Plot("./testmacd.png")
	e3 := time.Since(s)

	log.Printf("Timings:\n\tHistory: %s\n\tTimeSeries: %s\n\tMACD: %s", e1,e2,e3)
}

/**
 * In the following examples we fetch historical data between the desired time
 * stamps in the desired currency. Note poloniex doesn't give a granularity -
 * it just spews all the data it has.
 */
func polHist() {
	p := pol.Poloniex{false, []string{"USDT_ETH"}}
	recsP := p.Historic("USDT_ETH", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now())
	p.CSV("./outP.csv", recsP)
}

func gdaxHist() {
	g := gdax.GDAX{[]string{"ETH-USD"}}
	recs := g.Historic("ETH-USD", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 200)
	g.CSV("./outG.csv", recs)
}

/**
 * Examples of opening the websocket feeds to view realtime data from the exchanges.
 * Poloniex may be down.
 */
func polLive() {
	p := pol.Poloniex{true, []string{"USDT_ETH"}}
	go p.Live()
	time.Sleep(10 * time.Second)
	p.Enabled = false
}

func gdaxLive() {
	g := gdax.GDAX{[]string{"ETH-USD"}}

	// Asynchronously fetch data to messages channel.
	messages := make(chan gdax.WebsocketMatch)
	quit := make(chan bool)
	go g.Live(messages, quit)

	// Kill the livefeed after 10 seconds.
	go func() {
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	// Loop until something stops the socket feed (error or disabled)
	for msg := range messages {
		log.Println(msg)
	}
}
