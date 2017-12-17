# gdax_watcher

## TODO:
- Feed live data from each exchange into shared channel. Can have multiple exchange routines running at once, feed to consumers.
- Variety of consumers to make, start with MACD. Can generate plots per exchange, display the on top one another, or average them together by some means. Experiment.
- For now plots can be generated using go-charts, saved to a PNG and served to a frontend webpage. 
 - Later setup a websocket feed that can feed to a D3.JS frontend, allowing more interaction with the charts. You'll want to consider preprocessing as much data as you can rather than sending it in raw. 

## Keep It Simple: Start with 1 Exchange MACD.
### https://www.investopedia.com/terms/m/macd.asp
### For example EMA: https://github.com/AbenezerMamo/crypto-signal
### Sms alerts: https://www.twilio.com/sms
### Plots: https://github.com/wcharczuk/go-chart

1. For each day from today - 30 days:
2.  Day_MACD = 26_EMA - 12_EMA
3.  Day_Signal = 9_EMA
4.  Add to plot (Day = x, y1 & y2 = comps)

## Others: 
Keep an eye on this topic for more inspiration: https://github.com/topics/trading-bot