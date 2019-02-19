package data

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type MultisendResource struct {
	Coin  string `json:"coin"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (MultisendResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.TransactionOutput)

	return MultisendResource{
		Coin: data.Coin.Symbol,
		To: data.ToAddress.GetAddress(),
		Value: helpers.PipStr2Bip(data.Value),
	}
}
