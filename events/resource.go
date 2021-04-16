package events

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type BanResource struct {
	BlockID   uint64 `json:"height"`
	Timestamp string `json:"timestamp"`
}

func (BanResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	ban := model.(models.ValidatorBan)

	return BanResource{
		BlockID:   ban.BlockId,
		Timestamp: ban.Block.CreatedAt.Format(time.RFC3339),
	}
}

type AddressBanResource struct {
	BlockID   uint64             `json:"height"`
	Timestamp string             `json:"timestamp"`
	Validator resource.Interface `json:"validator"`
}

func (AddressBanResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	ban := model.(models.ValidatorBan)

	return AddressBanResource{
		BlockID:   ban.BlockId,
		Timestamp: ban.Block.CreatedAt.Format(time.RFC3339),
		Validator: new(validator.Resource).Transform(*ban.Validator),
	}
}
