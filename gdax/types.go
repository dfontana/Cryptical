package gdax

import (
	"encoding/json"
)

type GDAX struct {
	Enabled	bool
	Currencies []string
}

// WebsocketSubscribe takes in subscription information
type WebsocketSubscribe struct {
	Type      string `json:"type"`
	ProductID string `json:"product_id"`
}

// WebsocketOpen collates open orders
type WebsocketOpen struct {
	Type          string  `json:"type"`
	Time          string  `json:"time"`
	Sequence      int     `json:"sequence"`
	OrderID       string  `json:"order_id"`
	Price         float64 `json:"price,string"`
	RemainingSize float64 `json:"remaining_size,string"`
	Side          string  `json:"side"`
}

// WebsocketMatch holds match information
type WebsocketMatch struct {
	Type         string  `json:"type"`
	TradeID      int     `json:"trade_id"`
	Sequence     int     `json:"sequence"`
	MakerOrderID string  `json:"maker_order_id"`
	TakerOrderID string  `json:"taker_order_id"`
	Time         string  `json:"time"`
	ProductID		 string	 `json:"product_id"`
	Size         float64 `json:"size,string"`
	Price        float64 `json:"price,string"`
	Side         string  `json:"side"`
}

// WebsocketDone holds finished order information
type WebsocketDone struct {
	Type          string  `json:"type"`
	Time          string  `json:"time"`
	Sequence      int     `json:"sequence"`
	Price         float64 `json:"price,string"`
	OrderID       string  `json:"order_id"`
	Reason        string  `json:"reason"`
	Side          string  `json:"side"`
	RemainingSize float64 `json:"remaining_size,string"`
}

// History holds historic rate information
type Record struct {
	Time   int64 		`json:"time"`
	Low    float64	`json:"low"`
	High   float64	`json:"high"`
	Open   float64	`json:"open"`
	Close  float64	`json:"close"`
	Volume float64	`json:"volume"`
}

// UnmarshalJSON handles decomposing the returned array from GDAX into a series of record structs.
func (n *Record) UnmarshalJSON(buf []byte) error {
	tmp := []interface{}{&n.Time, &n.Low, &n.High, &n.Open, &n.Close, &n.Volume}
	if err := json.Unmarshal(buf, &tmp); err != nil {
		return err
	}
	return nil
}
