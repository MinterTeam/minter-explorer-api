package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type BuySwapPool struct {
	CoinToBuy          Coin   `json:"coin_to_buy"`
	CoinToSell         Coin   `json:"coin_to_sell"`
	ValueToBuy         string `json:"value_to_buy"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuySwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.BuySwapPoolData)

	return BuySwapPool{
		CoinToSell:         new(Coin).Transform(data.CoinToSell),
		CoinToBuy:          new(Coin).Transform(data.CoinToBuy),
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}

type SellSwapPool struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	ValueToSell       string `json:"value_to_sell"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellSwapPoolData)

	return SellSwapPool{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

type SellAllSwapPool struct {
	CoinToSell        Coin   `json:"coin_to_sell"`
	CoinToBuy         Coin   `json:"coin_to_buy"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellAllSwapPoolData)

	return SellAllSwapPool{
		CoinToSell:        new(Coin).Transform(data.CoinToSell),
		CoinToBuy:         new(Coin).Transform(data.CoinToBuy),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}
