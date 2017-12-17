package poloniex

type Poloniex struct {
	Enabled	bool
	Currencies []string
}

type Record struct {
	GlobalTradeID int64   `json:"globalTradeID"`
	TradeID       int64   `json:"tradeID"`
	Date          string  `json:"date"`
	Type          string  `json:"type"`
	Rate          float64 `json:"rate,string"`
	Amount        float64 `json:"amount,string"`
	Total 				float64 `json:"total,string"`
}