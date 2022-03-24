package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"time"
)

type Resource struct {
	Coin        resource.Interface `json:"coin"`
	Address     string             `json:"address"`
	Value       string             `json:"value"`
	Validator   resource.Interface `json:"validator"`
	ToValidator resource.Interface `json:"to_validator,omitempty"`
	BlockID     uint               `json:"height"`
	CreatedAt   string             `json:"created_at"`
	Type        string             `json:"type"`
}

const (
	typeUnbond    = "unbond"
	typeMoveStake = "move_stake"
)

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	unbond := model.(UnbondMoveStake)

	if unbond.ToValidator == nil {
		return Resource{
			Coin:      new(coins.IdResource).Transform(*unbond.Coin),
			Address:   unbond.Address.GetAddress(),
			Value:     helpers.PipStr2Bip(unbond.Value),
			Validator: new(validator.Resource).Transform(*unbond.FromValidator),
			BlockID:   unbond.BlockId,
			CreatedAt: unbond.CreatedAt.Format(time.RFC3339),
			Type:      typeUnbond,
		}
	}

	return Resource{
		Coin:        new(coins.IdResource).Transform(*unbond.Coin),
		Address:     unbond.Address.GetAddress(),
		Value:       helpers.PipStr2Bip(unbond.Value),
		Validator:   new(validator.Resource).Transform(*unbond.FromValidator),
		ToValidator: new(validator.Resource).Transform(*unbond.ToValidator),
		BlockID:     unbond.BlockId,
		CreatedAt:   unbond.CreatedAt.Format(time.RFC3339),
		Type:        typeMoveStake,
	}
}
