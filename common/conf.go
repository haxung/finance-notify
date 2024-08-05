package common

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
)

type ENV struct {
	Timeout  int      `toml:"timeout" json:"timeout"`
	Duration int      `toml:"duration" json:"duration"`
	RateURL  string   `toml:"rate_url" json:"rate_url"`
	RateIds  []string `toml:"rate_ids" json:"rate_ids"`
	CoinURL  string   `toml:"coin_url" json:"coin_url"`
	CoinIds  []string `toml:"coin_ids" json:"coin_ids"`
	Email    Email    `toml:"email" json:"email"`
	Rate     Rate     `toml:"rate" json:"rate"`
	Coin     Coin     `toml:"coin" json:"coin"`
	LogConf  LogConf  `toml:"log_conf" json:"log_conf"`
}

type Rate map[string][]float64
type Coin map[string][]float64

type Email struct {
	Server   string   `json:"server"`
	Port     int      `json:"port"`
	Username string   `json:"username"`
	Password string   `json:"password"`
	From     string   `json:"from"`
	To       []string `json:"to"`
}

var (
	env     = "env.toml"
	EnvConf = &ENV{}
)

func init() {
	flag.StringVar(&env, "c", env, "config env path")

	err := parseENV()
	if err != nil {
		panic("parseENV failed: " + err.Error())
	}
}

func parseENV() error {
	flag.Parsed()
	f, err := os.ReadFile(env)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(f, EnvConf)
	if err != nil {
		return err
	}

	return nil
}
