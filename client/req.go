package client

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type CoinData struct {
	Data []*RawCoinResp `json:"data"`
}

type RawCoinResp struct {
	Id                string `json:"id"`
	Rank              string `json:"rank"`
	Symbol            string `json:"symbol"`
	Name              string `json:"name"`
	Supply            string `json:"supply"`
	MaxSupply         string `json:"maxSupply"`
	MarketCapUsd      string `json:"marketCapUsd"`
	VolumeUsd24Hr     string `json:"volumeUsd24Hr"`
	PriceUsd          string `json:"priceUsd"`
	ChangePercent24Hr string `json:"changePercent24Hr"`
	Vwap24Hr          string `json:"vwap24Hr"`
}

type CoinResp struct {
	Id                string  `json:"id"`
	Rank              string  `json:"rand"`
	Symbol            string  `json:"symbol"`
	Vwap24Hr          float64 `json:"vwap24Hr"`
	PriceUsd          float64 `json:"priceUsd"`
	ChangePercent24Hr float64 `json:"changePercent24Hr"`
}

type CLI interface {
	GetCoinPrices(ids string) ([]*CoinResp, error)
	// GetRatePrices() error
	Close()
}

type cli struct {
	*http.Client
	rateURL string
	coinURL string
}

func NewCli(timeout int, rateURL, coinURL string) CLI {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	c := &http.Client{
		Transport: t,
		Timeout:   time.Duration(timeout) * time.Millisecond,
	}

	return &cli{
		Client:  c,
		rateURL: rateURL,
		coinURL: coinURL,
	}
}

func change(rawCoinResp []*RawCoinResp, res []*CoinResp) {
	for i, v := range rawCoinResp {
		price, _ := strconv.ParseFloat(v.PriceUsd, 64)
		vwap, _ := strconv.ParseFloat(v.Vwap24Hr, 64)
		changePercent, _ := strconv.ParseFloat(v.ChangePercent24Hr, 64)

		res[i] = &CoinResp{
			Id:                v.Id,
			Rank:              v.Rank,
			Symbol:            v.Symbol,
			Vwap24Hr:          vwap,
			PriceUsd:          price,
			ChangePercent24Hr: changePercent,
		}
	}
}

func (c *cli) GetCoinPrices(ids string) ([]*CoinResp, error) {
	resp, err := c.Get(c.coinURL + "?ids=" + ids)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	d := &CoinData{}
	err = json.Unmarshal(r, d)
	if err != nil {
		return nil, err
	}

	res := make([]*CoinResp, len(d.Data))
	change(d.Data, res)

	return res, nil
}

func (c *cli) Close() {
	c.Client.CloseIdleConnections()
}
