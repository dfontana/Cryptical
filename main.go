package main

import (
	pol "./poloniex"
	gdax "./gdax"
	"time"
	"log"
)

func main() {
	gdaxLive();
}

func polHist(){
	p := pol.Poloniex{false, []string{"USDT_ETH"}}
	recsP := p.Historic("USDT_ETH", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now())
	p.CSV("./outP.csv", recsP)
}

func gdaxHist(){
	g := gdax.GDAX{false, []string{"ETH-USD"}}
	recs := g.Historic("ETH-USD", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 200)
	g.CSV("./outG.csv", recs)
}

func polLive(){
	p := pol.Poloniex{true, []string{"USDT_ETH"}}
	go p.Live()
	time.Sleep(10 * time.Second)
	p.Enabled = false
}

func gdaxLive(){
	g := gdax.GDAX{true, []string{"ETH-USD"}}

	// Asynchronously fetch data to messages channel.
	messages := make(chan gdax.WebsocketMatch)
	quit := make(chan bool)
	go g.Live(messages, quit)

	// Kill the livefeed after 10 seconds.
	go func(){
		time.Sleep(10 * time.Second)
		quit <- true
	}()

	// Loop until something stops the socket feed (error or disabled)
	for msg := range messages{
		log.Println(msg)
	}
}