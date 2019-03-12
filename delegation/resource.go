package delegation

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	Coin   string `json:"coin"`
	PubKey string `json:"pub_key"`
	Value  string `json:"value"`
}

func (resource Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:   stake.Coin.Symbol,
		PubKey: stake.Validator.GetPublicKey(),
		Value:  helpers.PipStr2Bip(stake.Value),
	}
}
