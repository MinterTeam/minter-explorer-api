package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SellAllCoin struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellAllCoinData)

	return SellAllCoin{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
