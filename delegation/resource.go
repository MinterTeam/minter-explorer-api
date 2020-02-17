package delegation

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	validatorMeta "github.com/MinterTeam/minter-explorer-api/validator/meta"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type Resource struct {
	Coin          string             `json:"coin"`
	Value         string             `json:"value"`
	BipValue      string             `json:"bip_value"`
	PubKey        string             `json:"pub_key"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:          stake.Coin.Symbol,
		PubKey:        stake.Validator.GetPublicKey(),
		Value:         helpers.PipStr2Bip(stake.Value),
		BipValue:      helpers.PipStr2Bip(stake.BipValue),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*stake.Validator),
	}
}
