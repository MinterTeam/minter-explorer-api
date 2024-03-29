package helpers

import (
	"math"
	"math/big"
	"time"
)

// default amount of pips in 1 bip
var pipInBip = big.NewFloat(1000000000000000000)

var feeDefaultMultiplier = big.NewInt(100000000000000000)

// default amount of unit in one bip
const unitInBip = 10

func PipStr2Bip(value string) string {
	if value == "" {
		value = "0"
	}

	floatValue, err := new(big.Float).SetPrec(500).SetString(value)
	CheckErrBool(err)

	return new(big.Float).SetPrec(500).Quo(floatValue, pipInBip).Text('f', 18)
}

func Pip2Bip(pip *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(pip), pipInBip)
}

func Pip2BipStr(pip *big.Int) string {
	return Bip2Str(Pip2Bip(pip))
}

func Bip2Str(bip *big.Float) string {
	return bip.Text('f', 18)
}

func Bip2Pip(bip *big.Float) *big.Int {
	value, _ := bip.Mul(bip, pipInBip).Int(nil)
	return value
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
	if len(string) == 0 {
		return big.NewInt(0)
	}

	bInt, err := new(big.Int).SetString(string, 10)
	CheckErrBool(err)

	return bInt
}

func StrToBigFloat(string string) *big.Float {
	if len(string) == 0 {
		return big.NewFloat(0)
	}

	bFloat, err := new(big.Float).SetString(string)
	CheckErrBool(err)

	return bFloat
}
