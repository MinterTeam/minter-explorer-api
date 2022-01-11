package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type BuyCoin struct {
	CoinToBuy          Coin   `json:"coin_to_buy"`
	CoinToSell         Coin   `json:"coin_to_sell"`
	ValueToBuy         string `json:"value_to_buy"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.BuyCoinData)

	return BuyCoin{
		CoinToBuy:          new(Coin).Transform(data.CoinToBuy),
		CoinToSell:         new(Coin).Transform(data.CoinToSell),
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}
