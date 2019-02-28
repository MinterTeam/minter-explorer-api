package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/stake"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Resource struct {
	Status         *uint8               `json:"status"`
	Stake          *string              `json:"stake"`
	Part           *string              `json:"part"`
	DelegatorCount int                  `json:"delegator_count"`
	DelegatorList  []resource.Interface `json:"delegator_list"`
}

// Required extra params: activeValidatorIds and totalStake.
// activeValidatorIds: []uint64 active validator ids
// totalStake: string total stake of current active validator ids (by last block)
func (r Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	validator := model.(*models.Validator)
	delegators := resource.TransformCollection(validator.Stakes, stake.Resource{})
	activeValidators := params[0].([]uint64)
	totalStake := params[1].(string)

	var part string
	if helpers.InArray(validator.ID, activeValidators) && validator.TotalStake != nil {
		part = helpers.CalculatePercent(*validator.TotalStake, totalStake)
	}

	var stake string
	if validator.TotalStake != nil {
		stake = helpers.PipStr2Bip(*validator.TotalStake)
	}

	return Resource{
		Status:         validator.Status,
		Stake:          &stake,
		Part:           &part,
		DelegatorCount: len(delegators),
		DelegatorList:  delegators,
	}
}
