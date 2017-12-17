package main

import (
	pol "./poloniex"
	gdax "./gdax"
	"time"
)

func main() {
	p := pol.Poloniex{false, []string{"USDT_ETH"}}
	recsP := p.Historic("USDT_ETH", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now())
	p.CSV("./outP.csv", recsP)

	g := gdax.GDAX{false, []string{"ETH-USD"}}
	recs := g.Historic("ETH-USD", time.Date(2017, time.December, 14, 0, 0, 0, 0, time.Local), time.Now(), 200)
	g.CSV("./outG.csv", recs)
}