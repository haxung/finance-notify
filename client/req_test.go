package client

import (
	"log"
	"testing"
)

func TestGetCoinPrices(t *testing.T) {
	coinURL := "https://api.coincap.io/v2/assets"
	// Create a new client with the test server URL
	c := NewCli(5000, "", coinURL)

	// Call GetCoinPrices with no IDs
	res, err := c.GetCoinPrices("bitcoin,ethereum,tether,binance-coin,solana,xrp,dogecoin")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range res {
		log.Printf("%#v\n", v)
	}
}
