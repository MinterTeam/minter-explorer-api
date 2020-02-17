package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type BuyCoin struct {
	CoinToBuy          string `json:"coin_to_buy"`
	CoinToSell         string `json:"coin_to_sell"`
	ValueToBuy         string `json:"value_to_buy"`
	ValueToSell        string `json:"value_to_sell"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.BuyCoinTxData)
	model := params[0].(models.Transaction)

	return BuyCoin{
		CoinToBuy:          data.CoinToBuy,
		CoinToSell:         data.CoinToSell,
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		ValueToSell:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}
