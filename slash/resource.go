package slash

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type Resource struct {
	BlockID   uint64             `json:"height"`
	Coin      resource.Interface `json:"coin"`
	Amount    string             `json:"amount"`
	Address   string             `json:"address"`
	Timestamp string             `json:"timestamp"`
	Validator resource.Interface `json:"validator,omitempty"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	slash := model.(models.Slash)

	slashResource := Resource{
		BlockID:   slash.BlockID,
		Coin:      new(coins.IdResource).Transform(*slash.Coin),
		Amount:    helpers.PipStr2Bip(slash.Amount),
		Address:   slash.Address.GetAddress(),
		Timestamp: slash.Block.CreatedAt.Format(time.RFC3339),
	}

	if slash.Validator != nil {
		slashResource.Validator = new(validator.Resource).Transform(*slash.Validator)
	}

	return slashResource
}
