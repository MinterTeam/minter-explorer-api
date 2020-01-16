package invalid_transaction

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"time"
)

type Resource struct {
	Hash      string `json:"hash"`
	Block     uint64 `json:"block"`
	Type      uint8  `json:"type"`
	From      string `json:"from"`
	Timestamp string `json:"timestamp"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.InvalidTransaction)

	return Resource{
		Hash:      tx.GetHash(),
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Type:      tx.Type,
		From:      tx.FromAddress.GetAddress(),
	}
}
