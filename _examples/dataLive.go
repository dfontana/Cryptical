package examples

import (
  "log"
  "time"

  gdaxClient "github.com/dfontana/Cryptical/gdax"
  poloClient "github.com/dfontana/Cryptical/poloniex"
)

/**
 * Examples of opening the websocket feeds to view realtime data from the exchanges.
 * Poloniex may be down.
 */
func polLive() {
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
