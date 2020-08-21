package invalid_transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type Resource struct {
	Hash      string `json:"hash"`
	Block     uint64 `json:"block"`
	Timestamp string `json:"timestamp"`
	Type      uint8  `json:"type"`
	From      string `json:"from"`
	Log       string `json:"log"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tx := model.(models.InvalidTransaction)

	var txData map[string]interface{}
	json.Unmarshal([]byte(tx.TxData), &txData)

	return Resource{
		Hash:      tx.GetHash(),
		Block:     tx.BlockID,
		Timestamp: tx.CreatedAt.Format(time.RFC3339),
		Type:      tx.Type,
		From:      tx.FromAddress.GetAddress(),
		Log:       txData["log"].(string),
	}
}
