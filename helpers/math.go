package helpers

import (
	"math/big"
	"strconv"
)

// default amount of pips in 1 bip
var pipInBip = big.NewFloat(1000000000000000000)

func PipStr2Bip(value string) string {
	if value == "" {
		return value
	}

	floatValue, err := new(big.Float).SetString(value)
	CheckErrBool(err)

	return new(big.Float).Quo(floatValue, pipInBip).String()
}

func Fee2Bip(value uint64) string {
	return PipStr2Bip(strconv.FormatUint(value * 1000000000000000, 10))
}