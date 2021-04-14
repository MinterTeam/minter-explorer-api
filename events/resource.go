package events

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type BanResource struct {
	BlockID uint64 `json:"height"`
}

func (BanResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	ban := model.(models.ValidatorBan)

	return BanResource{
		BlockID: ban.BlockId,
	}
}
