package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type AddLiquidity struct {
	Coin0          Coin   `json:"coin0"`
	Coin1          Coin   `json:"coin1"`
	Volume0        string `json:"volume0"`
	MaximumVolume1 string `json:"maximum_volume1"`
}

func (AddLiquidity) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.AddLiquidityData)

	return AddLiquidity{
		Coin0:          new(Coin).Transform(data.Coin0),
		Coin1:          new(Coin).Transform(data.Coin1),
		Volume0:        helpers.PipStr2Bip(data.Volume0),
		MaximumVolume1: helpers.PipStr2Bip(data.MaximumVolume1),
	}
}
