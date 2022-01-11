package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SellCoin struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	ValueToSell       string `json:"value_to_sell"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellCoinData)

	return SellCoin{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
