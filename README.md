# Cryptical

## TODO:
Data
- Feed live data from each exchange into shared channel. Can have multiple exchange routines running at once, feed to consumers.
- Variety of consumers to make, start with MACD. Can generate plots per exchange, display the on top one another, or average them together by some means. Experiment.
- For now plots can be generated using go-charts, saved to a PNG and served to a frontend webpage. 
 - Later setup a websocket feed that can feed to a D3.JS frontend, allowing more interaction with the charts. You'll want to consider preprocessing as much data as you can rather than sending it in raw. 

 Automation:
 - Add inference for MACD (or monitoring, give some thought.)
 - Manually compute Bollinger
 - Add inference/monitoring for Bollinger
 - Mix sources

 Plotting:
 - minMax can be shared
 - Saving images can be shared


## How to run:
See `main.go` for some commentary and examples. To try them out, swap in the function you want to run and call `go run main.go`. Easy peasy. Later these can be built into binaries for executing, if desired.

#### Poloniex
- Work out an MACD example for Poloniex historical. Can you get hours / days?
- Add lookup logic for currency pairs to poloniex

## Dependencies:
- github.com/gorilla/websocket
- github.com/preichenberger/go-gdax
- github.com/wcharczuk/go-chart

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot
