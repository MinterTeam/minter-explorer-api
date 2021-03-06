package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SetHaltBlock struct {
	Height uint64 `json:"height"`
	PubKey string `json:"pub_key"`
}

func (SetHaltBlock) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SetHaltBlockData)

	return SetHaltBlock{
		Height: data.GetHeight(),
		PubKey: data.GetPubKey(),
	}
}
