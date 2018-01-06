package poloniex

import (
  "time"
  "strconv"
)

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

type CandleStick struct {
  Date            date		 		 `json:"date"`
  High            float64      `json:"high"`
  Low             float64      `json:"low"`
  Open            float64      `json:"open"`
  Close           float64      `json:"close"`
  Volume          float64      `json:"volume"`
  QuoteVolume     float64      `json:"quoteVolume"`
  WeightedAverage float64      `json:"weightedAverage"`
}

type date struct {
  time.Time 
}

func (d *date) UnmarshalJSON(data []byte) error {
  i, err := strconv.ParseInt(string(data), 10, 64)
  if err != nil {
    return err
  }
  d.Time = time.Unix(i, 0)
  return nil
}