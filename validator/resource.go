package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator/meta"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	PublicKey      string             `json:"public_key"`
	Status         *uint8             `json:"status"`
	Meta           resource.Interface `json:"meta"`
	Stake          *string            `json:"stake"`
	Part           *string            `json:"part"`
	DelegatorCount int                `json:"delegator_count"`
}

type Params struct {
	TotalStake          string // total stake of current active validator ids (by last block)
	ActiveValidatorsIDs []uint64
}

// Required extra params: object type of Params.
func (r Resource) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)
	params := values[0].(Params)
	part, validatorStake := r.getValidatorPartAndStake(validator, params.TotalStake, params.ActiveValidatorsIDs)

	result := Resource{
		PublicKey:      validator.GetPublicKey(),
		Status:         validator.Status,
		Stake:          validatorStake,
		Part:           part,
		Meta:           new(meta.Resource).Transform(validator),
		DelegatorCount: len(validator.Stakes),
	}

	return result
}

// return validator stake and part of the total (%)
func (r Resource) getValidatorPartAndStake(validator models.Validator, totalStake string, validators []uint64) (*string, *string) {
	var part, stake *string

	if helpers.InArray(validator.ID, validators) && validator.TotalStake != nil {
		val := helpers.CalculatePercent(*validator.TotalStake, totalStake)
		part = &val
	}

	if validator.TotalStake != nil {
		val := helpers.PipStr2Bip(*validator.TotalStake)
		stake = &val
	}

	return part, stake
}
