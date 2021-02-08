package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type RemoveLiquidity struct {
	Coin0          Coin   `json:"coin0"`
	Coin1          Coin   `json:"coin1"`
	Liquidity      string `json:"liquidity"`
	MinimumVolume0 string `json:"minimum_volume0"`
	MinimumVolume1 string `json:"minimum_volume1"`
	Volume0        string `json:"volume0"`
	Volume1        string `json:"volume1"`
}

func (RemoveLiquidity) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data, model := txData.(*api_pb.RemoveLiquidityData), params[0].(models.Transaction)

	return RemoveLiquidity{
		Coin0:          new(Coin).Transform(data.Coin0),
		Coin1:          new(Coin).Transform(data.Coin1),
		Liquidity:      helpers.PipStr2Bip(data.Liquidity),
		MinimumVolume0: helpers.PipStr2Bip(data.MinimumVolume0),
		MinimumVolume1: helpers.PipStr2Bip(data.MinimumVolume1),
		Volume0:        helpers.PipStr2Bip(model.Tags["tx.volume0"]),
		Volume1:        helpers.PipStr2Bip(model.Tags["tx.volume1"]),
	}
}
