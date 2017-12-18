package main

import (
	"log"
	"time"

	"./common"
	gdax "./gdax"
	pol "./poloniex"
)

func main() {
	gdaxMACD()
}

func gdaxMACD() {
	g := gdax.GDAX{[]string{"ETH-USD"}}

	daysBack := 150

	// Past 150 days for ETH daily.
	records := g.Historic(g.Currencies[0], time.Now().AddDate(0, 0, -daysBack), time.Now(), 24*60*60)

	// Due to unreliability in gdax API, we have to check if more data was returned than requested.
	if len(records) > daysBack+1 {
		log.Fatal("GDAX API returned more records than asked for, invalidating MACD computation.")
	}

	// Reduce to array of close values & their times
	hist := make([]common.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = common.TimeSeries{
			val.Time,
			val.Close,
		}
	}

	// MACD: 12 fast, 26 slow, 9 signal
	comp := common.MACD{}
	if err := comp.Populate(hist, 12, 26, 9); err != nil {
		log.Fatal(err)
	}

	log.Println("MACD: ", len(comp.Entries))
}

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
