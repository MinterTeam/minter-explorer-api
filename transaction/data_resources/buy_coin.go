package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type BuyCoin struct {
	CoinToBuy          Coin   `json:"coin_to_buy"`
	CoinToSell         Coin   `json:"coin_to_sell"`
	ValueToBuy         string `json:"value_to_buy"`
	ValueToSell        string `json:"value_to_sell"`
	MaximumValueToSell string `json:"maximum_value_to_sell"`
}

func (BuyCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.BuyCoin)
	model := params[0].(models.Transaction)

	return BuyCoin{
		CoinToBuy:          new(Coin).Transform(data.CoinToBuy),
		CoinToSell:         new(Coin).Transform(data.CoinToSell),
		ValueToBuy:         helpers.PipStr2Bip(data.ValueToBuy),
		ValueToSell:        helpers.PipStr2Bip(model.Tags["tx.return"]),
		MaximumValueToSell: helpers.PipStr2Bip(data.MaximumValueToSell),
	}
}
