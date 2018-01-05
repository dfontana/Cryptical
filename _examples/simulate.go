package main

import (
	"log"
	"time"

	"github.com/dfontana/Cryptical"
	"github.com/dfontana/Cryptical/computation"
	"github.com/dfontana/Cryptical/poloniex"
)

// Just for the time being while testing
func main() {

	// The below code is the outline of a simulation one might write on historic data
	// to evaluate their trading model strategy. Here, we evaluate an MACD model.

	/** ============================================
	  =============== MODEL PREP STEP ================
	  ============================================ **/
	// Day in history I'd like to infer for, -1 day since Historic is inclusive)
	// So if I want to infer for the 21st, I must ask for up to the 20th's of data
	endDate := time.Date(2017, time.December, 20, 23, 59, 59, 0, time.Local)

	// These are in 5 min periods meaining we look back 12, 5, and 3.5 hours
	slow := 24 * 6
	fast := 10 * 6
	sign := 7 * 6
	gran := 300 // Seconds

	// Will probably want ~3 * Slow data points to ensure we have enough.
	startHist := endDate.Add(time.Duration(-3*slow*gran) * time.Second)
	records, err := poloniex.Historic("USDT_ETH", startHist, endDate, gran)
	if err != nil {
		log.Fatal(err)
	}

	// Convert said data to the needed format to run our MACD computation.
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Date.Time,
			val.High,
		}
	}

	// Compute the inital model and plot it for visual inspection
	comp := computation.MACD{
		Data: hist,
		Fast: fast,
		Slow: slow,
		Sign: sign,
	}
	if err := comp.Compute(); err != nil {
		log.Fatal(err)
	}
	comp.Plot("./inference0.png")

	/** ============================================
	=============== STRATEGY DEFINITION STEP =======
	============================================ **/

	pf := bot.Portfolio{0.7, 0.0}
	var lastAction computation.Trade
	strategy := func(m *computation.MACD) computation.Trade {
		// This is just an example and will perform really poorly
		// Look at last entry and take action
		val := m.Hist[len(m.Hist)-1]
		action := computation.Trade{
			computation.Hodl,
			0.0,
			0.0,
		}
		switch {
		case val > 5:
			// Dont repeat actions
			if lastAction.Type != computation.Sell {
				// Sell 75% of current portfolio
				sellAmt := 0.75 * pf.Crypto
				action = computation.Trade{
					computation.Sell,
					sellAmt,
					m.Data[len(m.Data)-1].Data,
				}
				lastAction = action
			}
		case val < -2:
			// Dont repeat actions
			if lastAction.Type != computation.Buy {
				usdVal := 0.75 * pf.USD // Buy with 75% your funds
				if usdVal < 0.3 {
					// Too small, bail.
					action = computation.Trade{
						computation.Hodl,
						0.0,
						0.0,
					}
				} else {
					action = computation.Trade{
						computation.Buy,
						usdVal / m.Data[len(m.Data)-1].Data,
						m.Data[len(m.Data)-1].Data,
					}
				}
				lastAction = action
			}
		}
		return action
	}
	comp.Strategy = strategy

	/** ============================================
	=============== SIMULATION STEP ================
	============================================ **/
	tf := bot.Timeframe{endDate, gran, 24 * time.Hour}
	bot.Simulate(&comp, &pf, tf)
}
