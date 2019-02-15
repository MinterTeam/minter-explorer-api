package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SellAllCoinDataResource struct {
	CoinToSell        string `json:"coin_to_sell"`
	CoinToBuy         string `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllCoinDataResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SellAllCoinData)

	return SellAllCoinDataResource{
		CoinToSell:        data.CoinToSell,
		CoinToBuy:         data.CoinToBuy,
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

func (resource SellAllCoinDataResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SellAllCoinData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
