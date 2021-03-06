package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SellCoin struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	ValueToSell       string `json:"value_to_sell"`
	ValueToBuy        string `json:"value_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellCoinData)
	model := params[0].(models.Transaction)

	return SellCoin{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		ValueToBuy:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
