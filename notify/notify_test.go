package notify

import (
	"log"
	"testing"

	"finance-notify/client"
)

func TestNotify(t *testing.T) {
	err := Notify("hello, world")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCoinResp(t *testing.T) {
	coinURL := "https://api.coincap.io/v2/assets"
	// Create a new client with the test server URL
	c := client.NewCli(15000, "", coinURL)

	// Call GetCoinPrices with no IDs
	res, err := c.GetCoinPrices("bitcoin,ethereum,tether,binance-coin,solana,xrp,dogecoin")
	if err != nil {
		log.Fatal(err)
	}

	println(CoinResp(res))
}
