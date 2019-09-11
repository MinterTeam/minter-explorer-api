package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/stake"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	Status         *uint8               `json:"status"`
	Name           *string              `json:"name"`
	Description    *string              `json:"description"`
	IconUrl        *string              `json:"icon_url"`
	SiteUrl        *string              `json:"site_url"`
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

	var part *string
	if helpers.InArray(validator.ID, activeValidators) && validator.TotalStake != nil {
		val := helpers.CalculatePercent(*validator.TotalStake, totalStake)
		part = &val
	}

	var stake *string
	if validator.TotalStake != nil {
		val := helpers.PipStr2Bip(*validator.TotalStake)
		stake = &val
	}

	return Resource{
		Status:         validator.Status,
		Name:           validator.Name,
		Description:    validator.Description,
		IconUrl:        validator.IconUrl,
		SiteUrl:        validator.SiteUrl,
		Stake:          stake,
		Part:           part,
		DelegatorCount: len(delegators),
		DelegatorList:  delegators,
	}
}
