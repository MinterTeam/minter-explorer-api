package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SellAllCoinResource struct {
	CoinToSell        string `json:"coin_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllCoinResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SellAllCoinData)

	return SellAllCoinResource{
		CoinToSell:        data.CoinToSell,
		CoinToBuy:         data.CoinToBuy,
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

func (resource SellAllCoinResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SellAllCoinData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
