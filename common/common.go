package common

import (
	"math"
)

func IsNotify(symbol string, curr, percent float64) bool {
	coin := EnvConf.Coin[symbol]
	if curr <= coin[0] || curr >= coin[1] || math.Abs(percent) >= coin[3] {
		return true
	}

	return false
}
