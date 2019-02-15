package invalid_transaction

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"time"
)

type Resource struct {
	Hash      string                 `json:"hash"`
	Block     uint64                 `json:"block"`
	Timestamp time.Time              `json:"timestamp"`
	Type      string                 `json:"type"`
	From      string                 `json:"from"`
}

func (Resource) Transform(model resource.ItemInterface) resource.Interface {
	tx := model.(models.InvalidTransaction)

	return Resource{
		Hash: tx.GetHash(),
		Block: tx.BlockID,
		Timestamp: tx.CreatedAt,
		Type: transaction.GetTypeAsText(tx.Type),
		From: tx.FromAddress.Address,
	}
}
