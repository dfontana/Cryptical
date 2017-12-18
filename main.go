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

	// Past 30 days for ETH daily.
	records := g.Historic(g.Currencies[0], time.Now().AddDate(0, 0, -30), time.Now(), 24*60*60)
	log.Println(len(records))

	// Reduce to array of close values
	hist := make([]float64, len(records))
	for i, val := range records {
		hist[i] = val.Close
	}

	// MACD: 12 fast, 26 slow, 9 signal
	macd, _, err := common.MACD(hist, 12, 26, 9)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MACD: ", len(macd))
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
