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
		CoinToBuy:          new(Coin).Transform(data.Coins[len(data.Coins)-1]),
		CoinToSell:         new(Coin).Transform(data.Coins[0]),
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}

type SellSwapPool struct {
	CoinToBuy         Coin   `json:"coin_to_buy"`
	CoinToSell        Coin   `json:"coin_to_sell"`
	ValueToSell       string `json:"value_to_sell"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

type LimitOrderData struct {
	OrderId       uint64 `json:"order_id"`
	Buy           string `json:"buy"`
	Sell          string `json:"sell"`
	SellerAddress string `json:"seller_address"`
}

func (SellSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellSwapPoolData)

	return SellSwapPool{
		CoinToSell:        new(Coin).Transform(data.Coins[0]),
		CoinToBuy:         new(Coin).Transform(data.Coins[len(data.Coins)-1]),
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

type SellAllSwapPool struct {
	CoinToBuy         Coin   `json:"coin_to_buy"`
	CoinToSell        Coin   `json:"coin_to_sell"`
	MinimumValueToBuy string `json:"minimum_value_to_buy"`
}

func (SellAllSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SellAllSwapPoolData)

	return SellAllSwapPool{
		CoinToSell:        new(Coin).Transform(data.Coins[0]),
		CoinToBuy:         new(Coin).Transform(data.Coins[len(data.Coins)-1]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

type CreateSwapPool struct {
	Coin0   Coin   `json:"coin0"`
	Coin1   Coin   `json:"coin1"`
	Volume0 string `json:"volume0"`
	Volume1 string `json:"volume1"`
}

func (CreateSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.CreateSwapPoolData)

	return CreateSwapPool{
		Coin0:   new(Coin).Transform(data.Coin0),
		Coin1:   new(Coin).Transform(data.Coin1),
		Volume0: helpers.PipStr2Bip(data.Volume0),
		Volume1: helpers.PipStr2Bip(data.Volume1),
	}
}

type AddLimitOrder struct {
	CoinToSell  Coin   `json:"coin_to_sell"`
	CoinToBuy   Coin   `json:"coin_to_buy"`
	ValueToBuy  string `json:"value_to_buy"`
	ValueToSell string `json:"value_to_sell"`
}

func (AddLimitOrder) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.AddLimitOrderData)

	return AddLimitOrder{
		CoinToBuy:   new(Coin).Transform(data.CoinToBuy),
		CoinToSell:  new(Coin).Transform(data.CoinToSell),
		ValueToBuy:  helpers.PipStr2Bip(data.ValueToBuy),
		ValueToSell: helpers.PipStr2Bip(data.ValueToSell),
	}
}

type RemoveLimitOrder struct {
	Id uint64 `json:"id"`
}

func (RemoveLimitOrder) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.RemoveLimitOrderData)
	return RemoveLimitOrder{
		Id: data.Id,
	}
}
