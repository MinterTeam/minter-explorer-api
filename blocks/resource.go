package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"time"
)

type Resource struct {
	ID          uint64              `json:"height"`
	Size        uint64              `json:"size"`
	NumTxs      uint32              `json:"txCount"`
	BlockTime   uint64              `json:"blockTime"`
	Timestamp   time.Time           `json:"timestamp"`
	BlockReward string              `json:"reward"`
	Hash        string              `json:"hash"`
	Validators  []*models.Validator `json:"validators"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	block := model.(models.Block)

	return Resource{
		ID:          block.ID,
		Size:        block.Size,
		NumTxs:      block.NumTxs,
		BlockTime:   uint64(block.BlockTime),
		Timestamp:   block.CreatedAt,
		BlockReward: helpers.PipStr2Bip(block.BlockReward),
		Hash:        block.Hash,
		Validators:  block.Validators,
	}
}
