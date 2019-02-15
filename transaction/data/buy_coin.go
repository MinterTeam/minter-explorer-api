package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type BuyCoinResource struct {
	CoinToBuy          string `json:"coin_to_buy"`
	ValueToBuy         string `json:"value_to_buy"`
	CoinToSell         string `json:"coin_to_sell"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoinResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.BuyCoinData)

	return BuyCoinResource{
		CoinToBuy:          data.CoinToBuy,
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		CoinToSell:         data.CoinToSell,
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}

func (resource BuyCoinResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.BuyCoinData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
