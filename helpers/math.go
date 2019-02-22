package helpers

import (
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

	floatValue, err := new(big.Float).SetString(value)
	CheckErrBool(err)

	return new(big.Float).Quo(floatValue, pipInBip).String()
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
