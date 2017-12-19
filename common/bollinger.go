package common

// Bollinger handles generating bollinger band plots from the given historical
// price data. Populate the struct, then call its methods.
type Bollinger struct {
	history	[]float64
}

// Plot will create a bollinger plot from data stored in the type,
// saved to the given path.
func (b *Bollinger) Plot(path string) {

}