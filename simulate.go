package main

import (
  "log"
  "fmt"
  "time"

  "github.com/dfontana/Cryptical/plot"
  poloClient "github.com/dfontana/Cryptical/poloniex"
)

// Just for the time being while testing
func main() {

  // The below code is the outline of a simulation one might write on historic data
  // to evaluate their trading model strategy. Here, we evaluate an MACD model.

  // Day in history I'd like to infer for, -1 day since Historic is inclusive)
  // So if I want to infer for the 21st, I must ask for up to the 20th's of data
  startDate := time.Date(2017, time.December, 20, 23, 59, 59, 0, time.Local)
  
	// These are in 5 min periods meaining we look back 12, 5, and 3.5 hours
	slow := 24 * 6
	fast := 10 * 6
	sign := 7 * 6
	gran := 300 // Seconds

	// Will probably want ~3 * Slow data points to ensure we have enough.
	startHist := startDate.Add(-3 * slow * gran * time.Second) 
  records, err := poloClient.Historic("USDT_ETH", startHist, startDate, gran)
  if err != nil {
    log.Fatal(err)
  }

  // Convert said data to the needed format to run our MACD computation.
  hist := make([]plot.TimeSeries, len(records))
  for i, val := range records {
    hist[i] = plot.TimeSeries{
      val.Date.Time,
      val.High,
    }
  }

  // Compute MACD
  comp := plot.MACD{}
  if err := comp.Compute(hist, fast, slow, sign); err != nil {
    log.Fatal(err)
  }

  // Output the plot, just for our inspection
  comp.Plot("./inference.png")

  // Prepare the input and output channels
  trades := make(chan plot.Trade)
  soldprices := make(chan float64)

  // Start inference, for when data begins to flow
  go comp.Inference(soldprices, trades)

  // Get the fake historical data for the day after our start date
  records, err = poloClient.Historic("USDT_ETH", startDate, startDate.Add(24 * time.Hour), gran)
  if err != nil {
    log.Fatal(err)
  }

  // Feed in data to channel, we'll use Highs for this.
  go func() {
    for _,record := records {
      soldprices <- record.High
    }
  }()

  // Recieve what the bot is doing, until it has closed the channel.
  // This would be the chance to keep a running total of what the bot has done
  // in a log, as well as determining its final profit. For now, we print.
  for trade := range trades {
    var typ string
    switch(trade.Type) {
      case "sell":
        typ = "Sold"
      case "buy":
        typ = "Bought"
      case "hold":
        typ = "Held"
    }
    log.Println(fmt.Sprintf("%s %.5f for $%.2f\n", smt, trade.Crypto, trade.USD))
  }
  log.Println("Done.")
}
