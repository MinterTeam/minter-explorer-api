package validators

import "github.com/MinterTeam/minter-explorer-api/v2/core"

// Get IDs of active validators
func getActiveValidatorIDs(explorer *core.Explorer) []uint {
	return explorer.Cache.Get("active_validators", func() interface{} {
		return explorer.ValidatorRepository.GetActiveValidatorIds()
	}, CacheBlocksCount).([]uint)
}

// Get total stake of active validators
func getTotalStakeByActiveValidators(explorer *core.Explorer, validators []uint) string {
	return explorer.Cache.Get("validators_total_stake", func() interface{} {
		if len(validators) == 0 {
			return "0"
		}

		return explorer.ValidatorRepository.GetTotalStakeByActiveValidators(validators)
	}, CacheBlocksCount).(string)
}
