package delegation

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Coin         resource.Interface `json:"coin"`
	Value        string             `json:"value"`
	BipValue     string             `json:"bip_value"`
	Validator    resource.Interface `json:"validator"`
	IsWaitlisted bool               `json:"is_waitlisted"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:         new(coins.IdResource).Transform(*stake.Coin),
		Value:        helpers.PipStr2Bip(stake.Value),
		BipValue:     helpers.PipStr2Bip(stake.BipValue),
		Validator:    new(validator.Resource).Transform(*stake.Validator),
		IsWaitlisted: stake.IsKicked,
	}
}
