package chart

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"time"
)

type TransactionResource struct {
	Date    string `json:"date"`
	TxCount uint64 `json:"transaction_count"`
}

func (TransactionResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := model.(transaction.TxCountChartData)

	return TransactionResource{
		Date:    data.Time.Format(time.RFC3339),
		TxCount: data.Count,
	}
}
