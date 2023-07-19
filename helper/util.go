package helper

import (
	"math"
	"math/big"
)

func ConvertBalanceToFloat(balance *big.Int, decimalPlaces int) *big.Float {
	scaleFactor := new(big.Float).SetInt(big.NewInt(int64(math.Pow10(decimalPlaces))))
	floatBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), scaleFactor)
	return floatBalance
}
