package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SellAllCoin struct {
	CoinToSell        string `json:"coin_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllCoin) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SellAllCoinTxData)

	return SellAllCoin{
		CoinToSell:        data.CoinToSell,
		CoinToBuy:         data.CoinToBuy,
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
