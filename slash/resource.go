package slash

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	validatorMeta "github.com/MinterTeam/minter-explorer-api/validator/meta"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"time"
)

type Resource struct {
	BlockID       uint64             `json:"block"`
	Coin          string             `json:"coin"`
	Amount        string             `json:"amount"`
	Address       string             `json:"address"`
	Validator     string             `json:"validator"`
	ValidatorMeta resource.Interface `json:"validator_meta"`
	Timestamp     string             `json:"timestamp"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	slash := model.(models.Slash)

	return Resource{
		BlockID:       slash.BlockID,
		Coin:          slash.Coin.Symbol,
		Amount:        helpers.PipStr2Bip(slash.Amount),
		Address:       slash.Address.GetAddress(),
		Validator:     slash.Validator.GetPublicKey(),
		Timestamp:     slash.Block.CreatedAt.Format(time.RFC3339),
		ValidatorMeta: new(validatorMeta.Resource).Transform(*slash.Validator),
	}
}
