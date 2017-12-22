package main

import (
	"log"
	"os"
	gdax "github.com/preichenberger/go-gdax"
	_ "github.com/joho/godotenv/autoload"
)

const (
	// Trailing API point
	TRAILING_API = "/users/self/trailing-volume"

	// Taker Fees
	VOLUME0 = 0.003
	VOLUME1 = 0.0024
	VOLUME25 = 0.0022
	VOLUME5 = 0.0019
	VOLUME10 = 0.0015
	VOLUME20 = 0.001

	// Maker fee is 0.
)

// BreakEven determines the minimum value of ETH you'd want to sell for to make
// back the fees you incur during the transaction. These fees are recorded in
// the fee constants. There are two return values: the first is if you pay a
// taker fee, the second would be a maker fee.
// currency: The currency pair you bought
// usd: The amount you bought that pair for w/ fee accounted for
func BreakEven(currency string, usd float64) (float64, float64) {
	// Make client
	secret := os.Getenv("GDAX_SECRET")
	key := os.Getenv("GDAX_KEY")
	passphrase := os.Getenv("GDAX_KEY")
	c := gdax.NewClient(secret, key, passphrase)

	// Declare trailing type (array of these structs)
	var trailings []struct{
		ID							string 	`json:"product_id"`
		ExchangeVolume 	float64	`json:"exchange_volume, string"`
		Volume					float64	`json:"volume,string"`
		RecordedAt 			string	`json:"recorded_at"`
	}

	// Fill in struct array
	if _, err := c.Request("GET", TRAILING_API, nil, &trailings); err != nil {
		log.Fatal(err)
	}
	
	// Find the volume % for currency
	var vol float64
	for _,item := range trailings {
		if item.ID == currency {
			vol = item.Volume / item.ExchangeVolume
			break
		}
	}

	// Get the fee amount
	var volumeFee float64
	switch {
		case 0.0 <= vol && vol < 1.0:
			volumeFee = VOLUME0
		case 1.0 <= vol && vol < 2.5:
			volumeFee = VOLUME1
		case 2.5 <= vol && vol < 5.0:
			volumeFee = VOLUME25
		case 5.0 <= vol && vol < 10.0:
			volumeFee = VOLUME5
		case 10.0 <= vol && vol < 20.0:
			volumeFee = VOLUME10
		case 20.0 <= vol:
			volumeFee = VOLUME20
	}

	// Compute amount needed to offset selling fee
	taker := usd / 1 + volumeFee
	maker := usd
	return taker, maker
}