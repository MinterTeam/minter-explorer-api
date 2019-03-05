package chart

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"time"
)

type TransactionResource struct {
	Date    string `json:"date"`
	TxCount uint64 `json:"txCount"`
}

func (TransactionResource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	data := model.(transaction.TxCountChartData)

	return TransactionResource{
		Date:    data.Time.Format(time.RFC3339),
		TxCount: data.Count,
	}
}
