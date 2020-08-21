package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SellAllCoin struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	ValueToSell       string `json:"value_to_sell"`
	ValueToBuy        string `json:"value_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellAllCoinData)
	model := params[0].(models.Transaction)

	return SellAllCoin{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		ValueToSell:       helpers.PipStr2Bip(model.Tags["tx.sell_amount"]),
		ValueToBuy:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
