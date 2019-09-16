package aggregated_reward

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	validatorMeta "github.com/MinterTeam/minter-explorer-api/validator/meta"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	FromBlockID   uint64             `json:"from_block_id"`
	ToBlockID     uint64             `json:"to_block_id"`
	Role          string             `json:"role"`
	Amount        string             `json:"amount"`
	Address       string             `json:"address"`
	Validator     string             `json:"validator"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	reward := model.(models.AggregatedReward)

	return Resource{
		FromBlockID:   reward.FromBlockID,
		ToBlockID:     reward.ToBlockID,
		Role:          reward.Role,
		Amount:        helpers.PipStr2Bip(reward.Amount),
		Address:       reward.Address.GetAddress(),
		Validator:     reward.Validator.GetPublicKey(),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*reward.Validator),
	}
}

