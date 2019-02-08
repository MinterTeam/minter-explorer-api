package helpers

import (
	"math/big"
)

func PipStr2Bip(value string) string {
	if value == "" {
		return value
	}

	// default amount pip in 1 bip
	pipInBip := big.NewFloat(1000000000000000000)
	floatValue, err := new(big.Float).SetString(value)
	CheckErrBool(err)

	return new(big.Float).Quo(floatValue, pipInBip).String()
}
