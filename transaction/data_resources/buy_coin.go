package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type BuyCoin struct {
	CoinToBuy          string `json:"coin_to_buy"`
	ValueToBuy         string `json:"value_to_buy"`
	CoinToSell         string `json:"coin_to_sell"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoin) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.BuyCoinTxData)

	return BuyCoin{
		CoinToBuy:          data.CoinToBuy,
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		CoinToSell:         data.CoinToSell,
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}
