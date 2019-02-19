package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Delegate struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Stake  string `json:"stake"`
}

func (Delegate) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.DelegateTxData)

	return Delegate{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Stake:  helpers.PipStr2Bip(data.Stake),
	}
}
