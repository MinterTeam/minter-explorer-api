package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"time"
)

type BlockResource struct {
	ID          uint64              `json:"height"`
	Size        uint64              `json:"size"`
	NumTxs      uint32              `json:"txCount"`
	BlockTime   uint64              `json:"blockTime"`
	Timestamp   time.Time           `json:"timestamp"`
	BlockReward string              `json:"reward"`
	Hash        string              `json:"hash"`
	Validators  []*models.Validator `json:"validators"`
}

func TransformBlock(model models.Block) BlockResource {
	return BlockResource{
		ID:          model.ID,
		Size:        model.Size,
		NumTxs:      model.NumTxs,
		BlockTime:   model.BlockTime,
		Timestamp:   model.CreatedAt,
		BlockReward: helpers.PipStr2Bip(model.BlockReward),
		Hash:        model.Hash,
		Validators:  model.Validators,
	}
}

func TransformBlockCollection(models []models.Block) []BlockResource {
	var data []BlockResource
	for _, item := range models {
		data = append(data, TransformBlock(item))
	}

	return data
}
