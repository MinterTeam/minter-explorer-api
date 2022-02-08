package helpers

import (
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math"
	"strconv"
	"strings"
)

func RemoveMinterPrefix(raw string) string {
	if len(raw) < 2 {
		return raw
	}

	return raw[2:]
}

const firstReward = 333
const blocksPerReward = 200000
const premineValue = 200 * 1000000 // in the mainnet the real value is 200 * 1000000
const startHeight = 5000001 + 4150000

func CalculateEmission(blockId uint64) uint64 {
	blockId += startHeight
	sum := uint64(premineValue)
	reward := firstReward
	high := int(math.Ceil(float64(blockId) / float64(blocksPerReward)))

	var i int
	for i = 1; i < high; i++ {
		sum += uint64(blocksPerReward * reward)
		reward -= 1
	}

	sum += (blockId % uint64(blocksPerReward)) * uint64(reward)

	return sum
}

func GetSymbolAndVersionFromStr(symbol string) (string, *uint64) {
	items := strings.Split(symbol, "-")
	baseSymbol := items[0]

	if baseSymbol == "LP" {
		return symbol, nil
	}

	if len(items) == 2 {
		version, _ := strconv.ParseUint(items[1], 10, 64)
		return baseSymbol, &version
	}

	return baseSymbol, nil
}

func GetSymbolAndDefaultVersionFromStr(symbol string) (string, uint64) {
	symbol, version := GetSymbolAndVersionFromStr(symbol)
	if version == nil {
		return symbol, 0
	}

	return symbol, *version
}

func GetCoinType(coinType models.CoinType) string {
	switch coinType {
	case models.CoinTypeBase:
		return "coin"
	case models.CoinTypeToken:
		return "token"
	case models.CoinTypePoolToken:
		return "pool_token"
	}

	return ""
}

func GetPoolIdFromToken(token string) uint {
	id, _ := strconv.ParseUint(token[2:], 10, 64)
	return uint(id)
}

func GetTokenContractAndChain(tokenContract *models.TokenContract) (contract string, chain string) {
	// for BIP and HUB
	if tokenContract.CoinId == 0 || tokenContract.CoinId == 1902 {
		return tokenContract.Eth, "minter"
	}

	contract = tokenContract.Bsc
	if len(tokenContract.Eth) != 0 {
		return tokenContract.Eth, "ethereum"
	}

	return contract, "bsc"
}
