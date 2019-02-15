package data

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type UnbondResource struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Value  string `json:"value"`
}

func (UnbondResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.UnbondData)

	return UnbondResource{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Value:  helpers.PipStr2Bip(data.Value),
	}
}
