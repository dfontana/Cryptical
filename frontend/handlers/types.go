package handlers

type Ping struct {
	Ping string
}

type MacdModelRequest struct {
	EndDate     string `json:"end_date" binding:"Required"`
	Fast        int    `json:"fast" binding:"Required"`
	Slow        int    `json:"slow" binding:"Required"`
	Signal      int    `json:"signal" binding:"Required"`
	Granularity int    `json:"granularity" binding:"Required"`
	Pair        string `json:"pair" binding:"Required"`
}
