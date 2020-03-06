package helpers

import (
	"math"
	"math/big"
	"time"
)

// default amount of pips in 1 bip
var pipInBip = big.NewFloat(1000000000000000000)

var feeDefaultMultiplier = big.NewInt(1000000000000000)

// default amount of unit in one bip
const unitInBip = 1000

func PipStr2Bip(value string) string {
	if value == "" {
		value = "0"
	}

	floatValue, err := new(big.Float).SetPrec(500).SetString(value)
	CheckErrBool(err)

	return new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Text('f', 18)
}

func Fee2Bip(value uint64) string {
	return PipStr2Bip(new(big.Int).Mul(feeDefaultMultiplier, new(big.Int).SetUint64(value)).String())
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
	return float64(nano) / float64(time.Second)
}

func Unit2Bip(units float64) float64 {
	return units / unitInBip
}

func Seconds2Nano(sec int) float64 {
	return float64(sec) * float64(time.Second)
}

func StringToBigInt(string string) *big.Int {
	bInt, err := new(big.Int).SetString(string, 10)
	CheckErrBool(err)

	return bInt
}
