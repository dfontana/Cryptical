package gdax

// WebsocketOpen collates open orders
type WsOpen struct {
  Type          string  `json:"type"`
  Time          string  `json:"time"`
  Sequence      int     `json:"sequence"`
  OrderID       string  `json:"order_id"`
  Price         float64 `json:"price,string"`
  RemainingSize float64 `json:"remaining_size,string"`
  Side          string  `json:"side"`
}

// WebsocketMatch holds match information
type WsMatch struct {
  Type         string  `json:"type"`
  TradeID      int     `json:"trade_id"`
  Sequence     int     `json:"sequence"`
  MakerOrderID string  `json:"maker_order_id"`
  TakerOrderID string  `json:"taker_order_id"`
  Time         string  `json:"time"`
  ProductID    string  `json:"product_id"`
  Size         float64 `json:"size,string"`
  Price        float64 `json:"price,string"`
  Side         string  `json:"side"`
}

// WebsocketDone holds finished order information
type WsDone struct {
  Type          string  `json:"type"`
  Time          string  `json:"time"`
  Sequence      int     `json:"sequence"`
  Price         float64 `json:"price,string"`
  OrderID       string  `json:"order_id"`
  Reason        string  `json:"reason"`
  Side          string  `json:"side"`
  RemainingSize float64 `json:"remaining_size,string"`
}