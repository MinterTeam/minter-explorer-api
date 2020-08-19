package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type Unbond struct {
	PubKey string `json:"pub_key"`
	Coin   Coin   `json:"coin"`
	Value  string `json:"value"`
}

func (Unbond) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.UnbondData)

	return Unbond{
		PubKey: data.PubKey,
		Coin:   new(Coin).Transform(data.Coin),
		Value:  helpers.PipStr2Bip(data.Value),
	}
}
