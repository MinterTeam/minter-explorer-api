package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type Multisend struct {
	List []Send `json:"list"`
}

func (Multisend) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.MultiSendTxData)

	list := make([]Send, len(data.List))
	for key, item := range data.List {
		list[key] = Send{}.Transform(&item).(Send)
	}

	return Multisend{list}
}

func (Multisend) TransformByTxOutput(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.TransactionOutput)

	return Send{
		Coin:  data.Coin.Symbol,
		To:    data.ToAddress.GetAddress(),
		Value: helpers.PipStr2Bip(data.Value),
	}
}
