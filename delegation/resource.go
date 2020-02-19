package delegation

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type Resource struct {
	Coin      string             `json:"coin"`
	Value     string             `json:"value"`
	BipValue  string             `json:"bip_value"`
	Validator resource.Interface `json:"validator"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:      stake.Coin.Symbol,
		Value:     helpers.PipStr2Bip(stake.Value),
		BipValue:  helpers.PipStr2Bip(stake.BipValue),
		Validator: new(validator.Resource).Transform(*stake.Validator),
	}
}
