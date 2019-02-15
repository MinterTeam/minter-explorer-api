package data

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SendResource struct {
	Coin  string `json:"coin"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (SendResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SendData)

	return SendResource{
		Coin:  data.Coin,
		To:    data.To,
		Value: helpers.PipStr2Bip(data.Value),
	}
}
