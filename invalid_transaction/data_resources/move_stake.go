package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type MoveStake struct {
	FromPubKey string `json:"from_pub_key"`
	ToPubKey   string `json:"to_pub_key"`
	Coin       Coin   `json:"coin"`
	Value      string `json:"value"`
}

func (MoveStake) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.MoveStakeData)

	return MoveStake{
		FromPubKey: data.FromPubKey,
		ToPubKey:   data.ToPubKey,
		Coin:       new(Coin).Transform(data.Coin),
		Value:      helpers.PipStr2Bip(data.Value),
	}
}
