package chart

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/transaction"
)

type TransactionResource struct {
	Date    string `json:"date"`
	TxCount uint64 `json:"txCount"`
}

func (TransactionResource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	data := model.(transaction.TxCountChartData)

	return TransactionResource{
		Date:    data.Time.Format("2006-01-02T15:04:05.999999-07:00"),
		TxCount: data.Count,
	}
}
