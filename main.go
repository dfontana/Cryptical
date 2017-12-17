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
	// If a fatal error occurs, a message is sent to errors 
	errors := make(chan error)
	messages := make(chan gdax.WebsocketMatch)
	go g.Live(messages, errors)

	// Kill the livefeed after 5 seconds.
	go func(){
		time.Sleep(5 * time.Second)
		g.Enabled = false
	}()

	// Loop until something stops the socket feed (error or disabled)
	for msg := range messages{
		log.Println(msg)
	}

	// If loop broke due to error, print it now
	err := <-errors
	if err != nil {
		log.Println(err)
	}
}