package stake

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Coin         resource.Interface `json:"coin"`
	Address      string             `json:"address"`
	Value        string             `json:"value"`
	BipValue     string             `json:"bip_value"`
	IsWaitlisted bool               `json:"is_waitlisted"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	stake := model.(models.Stake)

	return Resource{
		Coin:         new(coins.IdResource).Transform(*stake.Coin),
		Address:      stake.OwnerAddress.GetAddress(),
		Value:        helpers.PipStr2Bip(stake.Value),
		BipValue:     helpers.PipStr2Bip(stake.BipValue),
		IsWaitlisted: stake.IsKicked,
	}
}
