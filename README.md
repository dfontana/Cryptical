# gdax_watcher

## TODO:
Data
- Feed live data from each exchange into shared channel. Can have multiple exchange routines running at once, feed to consumers.
- Variety of consumers to make, start with MACD. Can generate plots per exchange, display the on top one another, or average them together by some means. Experiment.
- For now plots can be generated using go-charts, saved to a PNG and served to a frontend webpage. 
 - Later setup a websocket feed that can feed to a D3.JS frontend, allowing more interaction with the charts. You'll want to consider preprocessing as much data as you can rather than sending it in raw. 

 Plotting:
 - minMax can be shared
 - Saving images can be shared


## How to run:
See `main.go` for some commentary and examples. To try them out, swap in the function you want to run and call `go run main.go`. Easy peasy. Later these can be built into binaries for executing, if desired.

## Dependencies:
Install these by running the `go get` command for each one independently (ie `go get github.com/gammazero/nexus/client`)
Poloniex Websockets:
- github.com/gammazero/nexus/client
- github.com/gammazero/nexus/wamp

GDax Websockets:
- github.com/gorilla/websocket

Plotting:
- github.com/wcharczuk/go-chart

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot
