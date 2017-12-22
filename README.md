# Cryptical
[![GoDoc](https://godoc.org/github.com/dfontana/Cryptical?status.svg)](https://godoc.org/github.com/dfontana/Cryptical)

## TODO:
- Put main.go into multiple test files inside an "_examples" folder (gdax.go, poloniex.go)
- Test cases.

 Automation:
 - Add inference for MACD (or monitoring, give some thought.)
 - Add inference/monitoring for Bollinger
 - Mix sources

Data
- Feed live data from each exchange into shared channel. Can have multiple exchange routines running at once, feed to consumers.
- Variety of consumers to make, start with MACD. Can generate plots per exchange, display the on top one another, or average them together by some means. Experiment.
- For now plots can be generated using go-charts, saved to a PNG and served to a frontend webpage. 
 - Later setup a websocket feed that can feed to a D3.JS frontend, allowing more interaction with the charts. You'll want to consider preprocessing as much data as you can rather than sending it in raw. 


## How to run:
See `main.go` for some commentary and examples. To try them out, first you need to `go get github.com/dfontana/Cryptical`, then swap in the function you want to run and call `go run main.go`. Easy peasy.


## Dependencies:
- github.com/gorilla/websocket
- github.com/preichenberger/go-gdax
- github.com/wcharczuk/go-chart
- github.com/joho/godotenv

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot
