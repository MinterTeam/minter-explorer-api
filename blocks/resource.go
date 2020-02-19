package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"time"
)

type Resource struct {
	ID          uint64               `json:"height"`
	Size        uint64               `json:"size"`
	NumTxs      uint32               `json:"transaction_count"`
	BlockTime   float64              `json:"block_time"`
	Timestamp   string               `json:"timestamp"`
	BlockReward string               `json:"reward"`
	Hash        string               `json:"hash"`
	Validators  []resource.Interface `json:"validators"`
}

// lastBlockId - uint64 pointer to the last block height, optional field.
func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	block := model.(models.Block)

	return Resource{
		ID:          block.ID,
		Size:        block.Size,
		NumTxs:      block.NumTxs,
		BlockTime:   helpers.Nano2Seconds(block.BlockTime),
		Timestamp:   block.CreatedAt.Format(time.RFC3339),
		BlockReward: helpers.PipStr2Bip(block.BlockReward),
		Hash:        block.GetHash(),
		Validators:  resource.TransformCollection(block.BlockValidators, ValidatorResource{}),
	}
}

type ValidatorResource struct {
	Validator resource.Interface `json:"validator"`
	Signed    bool               `json:"signed"`
}

func (ValidatorResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	blockValidator := model.(models.BlockValidator)

	return ValidatorResource{
		Signed:    blockValidator.Signed,
		Validator: new(validator.Resource).Transform(blockValidator.Validator),
	}
}
