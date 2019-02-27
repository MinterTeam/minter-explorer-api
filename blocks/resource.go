package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"time"
)

type Resource struct {
	ID                uint64               `json:"height"`
	Size              uint64               `json:"size"`
	NumTxs            uint32               `json:"txCount"`
	BlockTime         uint64               `json:"blockTime"`
	Timestamp         time.Time            `json:"timestamp"`
	BlockReward       string               `json:"reward"`
	Hash              string               `json:"hash"`
	LatestBlockHeight *uint64              `json:"latestBlockHeight,omitempty"`
	Validators        []resource.Interface `json:"validators"`
}

// lastBlockId - uint64 pointer to the last block height, optional field.
func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	block := model.(models.Block)

	var lastBlockId *uint64
	if len(params) == 1 {
		lastBlockId = params[0].(*uint64)
	}

	return Resource{
		ID:                block.ID,
		Size:              block.Size,
		NumTxs:            block.NumTxs,
		BlockTime:         uint64(block.BlockTime),
		Timestamp:         block.CreatedAt,
		BlockReward:       helpers.PipStr2Bip(block.BlockReward),
		Hash:              block.GetHash(),
		LatestBlockHeight: lastBlockId,
		Validators:        resource.TransformCollection(block.BlockValidators, ValidatorResource{}),
	}
}

type ValidatorResource struct {
	PublicKey string `json:"publicKey"`
	Signed    *bool  `json:"signed"`
}

func (ValidatorResource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	blockValidator := model.(models.BlockValidator)

	return ValidatorResource{
		PublicKey: blockValidator.Validator.GetPublicKey(),
		Signed:    blockValidator.Signed,
	}
}
