# gdax_watcher

## TODO:
- Feed live data from each exchange into shared channel. Can have multiple exchange routines running at once, feed to consumers.
- Variety of consumers to make, start with MACD. Can generate plots per exchange, display the on top one another, or average them together by some means. Experiment.
- For now plots can be generated using go-charts, saved to a PNG and served to a frontend webpage. 
 - Later setup a websocket feed that can feed to a D3.JS frontend, allowing more interaction with the charts. You'll want to consider preprocessing as much data as you can rather than sending it in raw. 

#### Poloniex
- Work out an MACD example for Poloniex historical. Can you get hours / days?
- Finish cleaning up the Websocket feed.

## Dependencies:
Poloniex Websockets:
- github.com/gammazero/nexus/client
- github.com/gammazero/nexus/wamp

GDax Websockets:
- github.com/gorilla/websocket

Plotting:
- github.com/wcharczuk/go-chart

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot
