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