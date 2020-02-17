package stake

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type Resource struct {
	Coin     string `json:"coin"`
	Address  string `json:"address"`
	Value    string `json:"value"`
	BipValue string `json:"bip_value"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:     stake.Coin.Symbol,
		Address:  stake.OwnerAddress.GetAddress(),
		Value:    helpers.PipStr2Bip(stake.Value),
		BipValue: helpers.PipStr2Bip(stake.BipValue),
	}
}
