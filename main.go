package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"finance-notify/client"
	"finance-notify/common"
	"finance-notify/notify"
	"go.uber.org/zap"
)

func main() {
	err := common.InitLogger(&common.EnvConf.LogConf)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("serv start...")

	c := client.NewCli(common.EnvConf.Timeout, common.EnvConf.RateURL, common.EnvConf.CoinURL)
	defer c.Close()
	cc := make(chan struct{})
	go func() {
		listenCoin(c, cc)
	}()

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTERM)
	<-signals
	close(cc)
	log.Println("serv finished!!!")
	os.Exit(0)
}

func listenCoin(c client.CLI, cc chan struct{}) {
	ticker := time.NewTicker(time.Duration(20) * time.Second)
	ids := strings.Join(common.EnvConf.CoinIds, ",")
	for {
		select {
		case <-ticker.C:
			prices, err := c.GetCoinPrices(ids)
			if err != nil {
				zap.L().Error("GetCoinPrices error", zap.Error(err))
				continue
			}

			for _, price := range prices {
				if common.IsNotify(price.Symbol, price.PriceUsd, price.Vwap24Hr) {
					err = notify.Notify(notify.CoinResp(prices))
					if err != nil {
						zap.L().Error("Notify error", zap.Error(err))
					}
					break
				}
			}

		case <-cc:
			ticker.Stop()
			return
		}
	}
}
