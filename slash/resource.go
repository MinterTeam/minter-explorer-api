package slash

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"time"
)

type Resource struct {
	BlockID   uint64             `json:"block"`
	Coin      string             `json:"coin"`
	Amount    string             `json:"amount"`
	Address   string             `json:"address"`
	Timestamp string             `json:"timestamp"`
	Validator resource.Interface `json:"validator"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	slash := model.(models.Slash)

	return Resource{
		BlockID:   slash.BlockID,
		Coin:      slash.Coin.Symbol,
		Amount:    helpers.PipStr2Bip(slash.Amount),
		Address:   slash.Address.GetAddress(),
		Timestamp: slash.Block.CreatedAt.Format(time.RFC3339),
		Validator: new(validator.Resource).Transform(*slash.Validator),
	}
}
