package events

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type BanResource struct {
	FromBlockId   uint64 `json:"from_block_id"`
	FromTimestamp string `json:"from_timestamp"`
	ToBlockId     uint64 `json:"to_block_id"`
}

func (BanResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	ban := model.(models.ValidatorBan)

	return BanResource{
		FromBlockId:   ban.BlockId,
		FromTimestamp: ban.Block.CreatedAt.Format(time.RFC3339),
		ToBlockId:     ban.ToBlockId,
	}
}

type AddressBanResource struct {
	FromBlockId   uint64             `json:"from_block_id"`
	FromTimestamp string             `json:"from_timestamp"`
	ToBlockId     uint64             `json:"to_block_id"`
	Timestamp     string             `json:"timestamp"`
	Validator     resource.Interface `json:"validator"`
}

func (AddressBanResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	ban := model.(models.ValidatorBan)

	return AddressBanResource{
		FromBlockId:   ban.BlockId,
		FromTimestamp: ban.Block.CreatedAt.Format(time.RFC3339),
		ToBlockId:     ban.ToBlockId,
		Validator:     new(validator.Resource).Transform(*ban.Validator),
	}
}
