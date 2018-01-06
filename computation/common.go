package computation

import (
	"image/png"
	"os"

	"github.com/wcharczuk/go-chart"
)

// Sma is the Simple Moving Average for given slice
func Sma(closingPrices []TimeSeries) float64 {
	if len(closingPrices) == 0 {
		return 0
	}

	var sum float64
	for _, v := range closingPrices {
		sum += v.Data
	}

	return sum / float64(len(closingPrices))
}

// SaveImage renders the given graph and exports it as a PNG to the given path.
func SaveImage(graph chart.Chart, path string) error {
	// Write image to buffer
	collector := &chart.ImageWriter{}
	graph.Render(chart.PNG, collector)
	image, err := collector.Image()
	if err != nil {
		return err
	}

	// Save buffer to file (after encoding)
	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}
	png.Encode(outputFile, image)
	outputFile.Close()
	return nil
}

// MinMax returns the minimum and maximum value of a slice.
func MinMax(vals []float64) (float64, float64) {
	min := vals[0]
	max := vals[0]
	for _, val := range vals {
		if val < min {
			min = val
		}
		if val > max {
			max = val
		}
	}
	return min, max
}
