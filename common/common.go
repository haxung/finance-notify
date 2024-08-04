package common

import (
	"math"

	"go.uber.org/zap"
)

func IsNotify(symbol string, curr, percent float64) bool {
	coin := EnvConf.Coin[symbol]
	if len(coin) < 3 {
		zap.L().Error("coin data failed", zap.String("symbol", symbol), zap.Float64s("data", coin))
		return false
	}

	if curr <= coin[0] || curr >= coin[1] || math.Abs(percent) >= coin[2] {
		zap.L().Warn("coin notify", zap.String("symbol", symbol), zap.Float64s("data", coin))
		return true
	}

	return false
}
