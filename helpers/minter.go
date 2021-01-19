package helpers

import (
	"math"
	"strconv"
	"strings"
)

func RemoveMinterPrefix(raw string) string {
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
