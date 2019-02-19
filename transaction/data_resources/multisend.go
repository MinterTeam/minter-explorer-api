package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Multisend struct {
	Coin  string `json:"coin"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (Multisend) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.TransactionOutput)

	return Multisend{
		Coin: data.Coin.Symbol,
		To: data.ToAddress.GetAddress(),
		Value: helpers.PipStr2Bip(data.Value),
	}
}
