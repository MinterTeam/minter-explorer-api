package helpers

import (
	"math"
	"math/big"
	"strconv"
)

// default amount of pips in 1 bip
var pipInBip = big.NewFloat(1000000000000000000)

var feeDefaultMultiplier = uint64(1000000000000000)

func PipStr2Bip(value string) string {
	if value == "" {
		return value
	}

	floatValue, err := new(big.Float).SetPrec(500).SetString(value)
	CheckErrBool(err)

	return new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Text('f', 18)
}

func Fee2Bip(value uint64) string {
	return PipStr2Bip(strconv.FormatUint(value*feeDefaultMultiplier, 10))
}

func CalculatePercent(part string, total string) string {
	v1, err := new(big.Float).SetString(part)
	CheckErrBool(err)

	v2, err := new(big.Float).SetString(total)
	CheckErrBool(err)

	v1 = new(big.Float).Mul(v1, big.NewFloat(100))

	return new(big.Float).Quo(v1, v2).String()
}

func Round(value float64, precision int) float64 {
	return math.Round(value*math.Pow10(precision)) / math.Pow10(precision)
}

func Nano2Seconds(nano uint64) float64 {
	return float64(nano) / float64(1000000000)
}
