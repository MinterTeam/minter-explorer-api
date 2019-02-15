package data

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type DelegateResource struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Stake  string `json:"stake"`
}

func (DelegateResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.DelegateTxData)

	return DelegateResource{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Stake:  helpers.PipStr2Bip(data.Stake),
	}
}
