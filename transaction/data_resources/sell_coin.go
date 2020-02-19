package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type SellCoin struct {
	CoinToSell        string `json:"coin_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	ValueToSell       string `json:"value_to_sell"`
	ValueToBuy        string `json:"value_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.SellCoinTxData)
	model := params[0].(models.Transaction)

	return SellCoin{
		CoinToSell:        data.CoinToSell,
		CoinToBuy:         data.CoinToBuy,
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		ValueToBuy:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
