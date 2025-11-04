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
	"finance-notify/mail"

	"go.uber.org/zap"
)

func main() {
	err := common.InitLogger(&common.EnvConf.LogConf)
	if err != nil {
		log.Fatal(err)
	}

	zap.L().Warn("init logger success!!!")

	log.Println("serv start...")

	c := client.NewCli(common.EnvConf.Timeout, common.EnvConf.RateURL, common.EnvConf.CoinURL)
	defer c.Close()
	cc := make(chan struct{})
	go func() {
		listenCoin(c, cc)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTERM)
	<-signals
	close(cc)
	log.Println("serv finished!!!")
	os.Exit(0)
}

func listenCoin(c client.CLI, cc chan struct{}) {
	ticker := time.NewTicker(time.Duration(common.EnvConf.Duration) * time.Second)
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
				if common.IsNotify(price.Symbol, price.PriceUsd, price.ChangePercent24Hr) {
					err = mail.Notify(mail.CoinResp(prices))
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
