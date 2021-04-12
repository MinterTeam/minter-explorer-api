package pool

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"math/big"
)

func GetPoolChainFromStr(chainStr string) (chain []SwapChain) {
	json.Unmarshal([]byte(chainStr), &chain)
	return chain
}

func getVolumeInBip(price *big.Float, volume string) *big.Float {
	firstCoinBaseVolume := helpers.Pip2Bip(helpers.StringToBigInt(volume))
	return new(big.Float).Mul(firstCoinBaseVolume, price)
}
