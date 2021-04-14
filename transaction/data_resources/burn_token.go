package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type BurnToken struct {
	Coin  Coin   `json:"coin"`
	Value string `json:"value"`
}

func (BurnToken) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.BurnTokenData)

	return BurnToken{
		Coin:  new(Coin).Transform(data.Coin),
		Value: helpers.PipStr2Bip(data.Value),
	}
}
