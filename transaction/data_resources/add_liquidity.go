package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type AddLiquidity struct {
	Coin0          Coin   `json:"coin0"`
	Coin1          Coin   `json:"coin1"`
	PoolToken      Coin   `json:"pool_token"`
	Volume0        string `json:"volume0"`
	MaximumVolume1 string `json:"maximum_volume1"`
	Volume1        string `json:"volume1"`
	Liquidity      string `json:"liquidity"`
}

func (AddLiquidity) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.AddLiquidityData), params[0].(models.Transaction)

	return AddLiquidity{
		Coin0:          new(Coin).Transform(data.Coin0),
		Coin1:          new(Coin).Transform(data.Coin1),
		Volume0:        helpers.PipStr2Bip(data.Volume0),
		MaximumVolume1: helpers.PipStr2Bip(data.MaximumVolume1),
		Volume1:        helpers.PipStr2Bip(model.Tags["tx.volume1"]),
		Liquidity:      helpers.PipStr2Bip(model.Tags["tx.liquidity"]),
		PoolToken: new(Coin).Transform(&api_pb.Coin{
			Symbol: model.Tags["pool_token"],
			Id:     helpers.StrToUint64(model.Tags["pool_token_id"]),
		}),
	}
}
