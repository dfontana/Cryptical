
package examples

import (
  "log"
  "time"

  gdaxClient "github.com/dfontana/Cryptical/gdax"
  poloClient "github.com/dfontana/Cryptical/poloniex"
)

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
