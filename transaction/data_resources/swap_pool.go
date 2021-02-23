package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type BuySwapPool struct {
	Coins              []PoolCoin `json:"coins"`
	ValueToBuy         string     `json:"value_to_buy"`
	ValueToSell        string     `json:"value_to_sell"`
	MaximumValueToSell string     `json:"maximum_value_to_sell"`
}

func (BuySwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.BuySwapPoolData), params[0].(models.Transaction)

	return BuySwapPool{
		Coins:              new(PoolCoin).TransformCollection(data.Coins, model),
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		ValueToSell:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}

type SellSwapPool struct {
	Coins             []PoolCoin `json:"coins"`
	ValueToSell       string     `json:"value_to_sell"`
	ValueToBuy        string     `json:"value_to_buy"`
	MinimumValueToBuy string     `json:"minimum_value_to_buy"`
}

func (SellSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.SellSwapPoolData), params[0].(models.Transaction)

	return SellSwapPool{
		Coins:             new(PoolCoin).TransformCollection(data.Coins, model),
		ValueToSell:       helpers.PipStr2Bip(data.ValueToSell),
		ValueToBuy:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

type SellAllSwapPool struct {
	Coins             []PoolCoin `json:"coins"`
	ValueToSell       string     `json:"value_to_sell"`
	ValueToBuy        string     `json:"value_to_buy"`
	MinimumValueToBuy string     `json:"minimum_value_to_buy"`
}

func (SellAllSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.SellAllSwapPoolData), params[0].(models.Transaction)

	return SellAllSwapPool{
		Coins:             new(PoolCoin).TransformCollection(data.Coins, model),
		ValueToSell:       helpers.PipStr2Bip(model.Tags["tx.sell_amount"]),
		ValueToBuy:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MinimumValueToBuy: helpers.PipStr2Bip(data.MinimumValueToBuy),
	}
}

type CreateSwapPool struct {
	Coin0     Coin   `json:"coin0"`
	Coin1     Coin   `json:"coin1"`
	PoolToken Coin   `json:"pool_token"`
	Volume0   string `json:"volume0"`
	Volume1   string `json:"volume1"`
	Liquidity string `json:"liquidity"`
}

func (CreateSwapPool) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.CreateSwapPoolData), params[0].(models.Transaction)

	return CreateSwapPool{
		Coin0:     new(Coin).Transform(data.Coin0),
		Coin1:     new(Coin).Transform(data.Coin1),
		Volume0:   helpers.PipStr2Bip(data.Volume0),
		Volume1:   helpers.PipStr2Bip(data.Volume1),
		Liquidity: helpers.PipStr2Bip(model.Tags["tx.liquidity"]),
		PoolToken: Coin{
			ID:     helpers.StrToUint64(model.Tags["tx.pool_token_id"]),
			Symbol: model.Tags["tx.pool_token"],
		},
	}
}
