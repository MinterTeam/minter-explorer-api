package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type ResourceDetailed struct {
	PublicKey      string  `json:"public_key"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	IconUrl        *string `json:"icon_url"`
	SiteUrl        *string `json:"site_url"`
	Status         *uint8  `json:"status"`
	Stake          string  `json:"stake"`
	Commission     *uint64 `json:"commission"`
	Part           string  `json:"part"`
	DelegatorCount int     `json:"delegator_count"`
	MinStake       string  `json:"min_stake"`
}

type Params struct {
	TotalStake          string // total stake of current active validator ids (by last block)
	MinStake            string
	ActiveValidatorsIDs []uint
}

// Required extra params: object type of Params.
func (r ResourceDetailed) Transform(model resource.ItemInterface, values ...resource.ParamInterface) resource.Interface {
	validator := model.(models.Validator)
	params := values[0].(Params)

	validatorStakePart, totalStake := "0", "0"
	if helpers.InArray(validator.ID, params.ActiveValidatorsIDs) && validator.TotalStake != nil {
		validatorStakePart = helpers.CalculatePercent(*validator.TotalStake, params.TotalStake)
	}

	if validator.TotalStake != nil {
		totalStake = *validator.TotalStake
	}

	result := ResourceDetailed{
		PublicKey:      validator.GetPublicKey(),
		Status:         validator.Status,
		Stake:          helpers.PipStr2Bip(totalStake),
		Name:           validator.Name,
		Description:    validator.Description,
		IconUrl:        validator.IconUrl,
		SiteUrl:        validator.SiteUrl,
		Commission:     validator.Commission,
		Part:           validatorStakePart,
		MinStake:       helpers.PipStr2Bip(params.MinStake),
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
