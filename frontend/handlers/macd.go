package handlers

import (
	"time"

	"github.com/dfontana/Cryptical/computation"
	"github.com/dfontana/Cryptical/poloniex"
	m "gopkg.in/macaron.v1"
)

// ModelMACD returns an MACD model from the given parameters.
func ModelMACD(ctx *m.Context, mmr MacdModelRequest) {
	endDate, err := time.Parse("2006-01-02T15:04:05.999Z", mmr.EndDate)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}

	// Fetch data from Poloniex in appropriate format
	startHist := endDate.Add(time.Duration(-3*mmr.Slow*mmr.Granularity) * time.Second)
	records, err := poloniex.Historic(mmr.Pair, startHist, endDate, mmr.Granularity)
	if err != nil {
		ctx.Error(400, err.Error())
		return
	}
	hist := make([]computation.TimeSeries, len(records))
	for i, val := range records {
		hist[i] = computation.TimeSeries{
			val.Date.Time,
			val.High,
		}
	}

	// Compute the inital model and plot it for visual inspection
	model := computation.MACD{
		Data: hist,
		Fast: mmr.Fast,
		Slow: mmr.Slow,
		Sign: mmr.Signal,
	}
	if err := model.Compute(); err != nil {
		ctx.Error(400, err.Error())
		return
	}

	var res []MacdModelResponse
	for i := range model.Time {
		res = append(res, MacdModelResponse{
			Time:   model.Time[i],
			MACD:   model.MACD[i],
			Signal: model.Signal[i],
			Hist:   model.Hist[i],
		})
	}

	ctx.JSON(200, res)
}

// SimulateMACD returns a simulation of the given MACD model from the given
// paramters.
func SimulateMACD(ctx *m.Context) {

}
