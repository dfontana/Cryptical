package poloniex

import "time"

type WSOrderbook struct {
	Pair    string
	Event   string
	TradeID int64
	Type    string
	Rate    float64
	Amount  float64
	Total   float64
	TS      time.Time
}

type TickerEntry struct {
	Last        float64 `json:",string"`
	Ask         float64 `json:"lowestAsk,string"`
	Bid         float64 `json:"highestBid,string"`
	Change      float64 `json:"percentChange,string"`
	BaseVolume  float64 `json:"baseVolume,string"`
	QuoteVolume float64 `json:"quoteVolume,string"`
	IsFrozen    int64   `json:"isFrozen,string"`
	ID          int64   `json:"id"`
}