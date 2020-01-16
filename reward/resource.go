package reward

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	validatorMeta "github.com/MinterTeam/minter-explorer-api/validator/meta"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"time"
)

type Resource struct {
	BlockID       uint64             `json:"block"`
	Role          string             `json:"role"`
	Amount        string             `json:"amount"`
	Address       string             `json:"address"`
	Timestamp     string             `json:"timestamp"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	reward := model.(models.Reward)

	return Resource{
		BlockID:       reward.BlockID,
		Role:          reward.Role,
		Amount:        helpers.PipStr2Bip(reward.Amount),
		Address:       reward.Address.GetAddress(),
		Timestamp:     reward.Block.CreatedAt.Format(time.RFC3339),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*reward.Validator),
	}
}
