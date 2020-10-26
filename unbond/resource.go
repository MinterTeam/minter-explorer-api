package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Coin      resource.Interface `json:"coin"`
	Address   string             `json:"address"`
	Value     string             `json:"value"`
	Validator resource.Interface `json:"validator"`
	BlockID   uint               `json:"height"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	unbond := model.(models.Unbond)

	return Resource{
		Coin:      new(coins.IdResource).Transform(*unbond.Coin),
		Address:   unbond.Address.GetAddress(),
		Value:     helpers.PipStr2Bip(unbond.Value),
		Validator: new(validator.Resource).Transform(*unbond.Validator),
		BlockID:   unbond.BlockId,
	}
}
