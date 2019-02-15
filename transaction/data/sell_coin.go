package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SellCoinResource struct {
	CoinToSell        string `json:"coin_to_sell"`
	ValueToSell       string `json:"value_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellCoinResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SellCoinData)

	return SellCoinResource{
		CoinToSell:        data.CoinToSell,
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		CoinToBuy:         data.CoinToSell,
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

func (resource SellCoinResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SellCoinData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
