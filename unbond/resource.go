package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Coin             resource.Interface `json:"coin"`
	Address          string             `json:"address"`
	Value            string             `json:"value"`
	Validator        resource.Interface `json:"validator"`
	CreatedAtBlockID uint64             `json:"created_at_block_id"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	unbond := model.(models.Unbond)

	return Resource{
		Coin:             new(coins.IdResource).Transform(*unbond.Coin),
		Address:          unbond.Address.GetAddress(),
		Value:            helpers.PipStr2Bip(unbond.Value),
		Validator:        new(validator.Resource).Transform(*unbond.Validator),
		CreatedAtBlockID: 0, // TODO: fix
	}
}
