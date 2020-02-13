package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type ResourceDetailed struct {
	PublicKey      string  `json:"public_key"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	IconUrl        *string `json:"icon_url"`
	SiteUrl        *string `json:"site_url"`
	Status         *uint8  `json:"status"`
	Stake          *string `json:"stake"`
	Part           *string `json:"part"`
	DelegatorCount int     `json:"delegator_count"`
}

type Params struct {
	TotalStake          string // total stake of current active validator ids (by last block)
	ActiveValidatorsIDs []uint64
}

// Required extra params: object type of Params.
func (r ResourceDetailed) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)
	params := values[0].(Params)
	part, validatorStake := r.getValidatorPartAndStake(validator, params.TotalStake, params.ActiveValidatorsIDs)

	result := ResourceDetailed{
		PublicKey:      validator.GetPublicKey(),
		Status:         validator.Status,
		Stake:          validatorStake,
		Name:           validator.Name,
		Description:    validator.Description,
		IconUrl:        validator.IconUrl,
		SiteUrl:        validator.SiteUrl,
		Part:           part,
		DelegatorCount: len(validator.Stakes),
	}

	return result
}

// return validator stake and part of the total (%)
func (r ResourceDetailed) getValidatorPartAndStake(validator models.Validator, totalStake string, validators []uint64) (*string, *string) {
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

type Resource struct {
	PublicKey   string  `json:"public_key"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IconUrl     *string `json:"icon_url"`
	SiteUrl     *string `json:"site_url"`
	Status      *uint8  `json:"status"`
}

func (r Resource) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)

	result := Resource{
		PublicKey:   validator.GetPublicKey(),
		Status:      validator.Status,
		Name:        validator.Name,
		Description: validator.Description,
		IconUrl:     validator.IconUrl,
		SiteUrl:     validator.SiteUrl,
	}

	return result
}
