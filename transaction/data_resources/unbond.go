package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type Unbond struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Value  string `json:"value"`
}

func (Unbond) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.UnbondTxData)

	return Unbond{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Value:  helpers.PipStr2Bip(data.Value),
	}
}
