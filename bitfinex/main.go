package main

import (
	"time"
	"fmt"
	"github.com/bitfinexcom/bitfinex-api-go/v1"
)

func main() {
	// save this for later client := bitfinex.NewClient().Auth("api-key", "api-secret")
	client := bitfinex.NewClient() // removing auth to see if it works

	bitches, err := client.History.Trades("tETHBTC", time.Date(2017, time.December, 20, 23, 59, 59, 0, time.Local), time.Now(), 99999, false)
	fmt.Println(err)
	fmt.Println(bitches)

	// info, err := client.Account.Info()

	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(info)
	// }
}