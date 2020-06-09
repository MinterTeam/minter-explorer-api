package helpers

import "math"

func RemoveMinterPrefix(raw string) string {
	return raw[2:]
}

const firstReward = 333
const blocksPerReward = 200000
const premineValue = 200 * 1000000 // in the mainnet the real value is 200 * 1000000
const startHeight = 5000001

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
