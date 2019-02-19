package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SellCoin struct {
	CoinToSell        string `json:"coin_to_sell"`
	ValueToSell       string `json:"value_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellCoin) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SellCoinTxData)

	return SellCoin{
		CoinToSell:        data.CoinToSell,
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		CoinToBuy:         data.CoinToSell,
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
