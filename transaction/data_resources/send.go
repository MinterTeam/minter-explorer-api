package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type Send struct {
	Coin  Coin   `json:"coin"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func (Send) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.SendData)

	return Send{
		Coin:  new(Coin).Transform(data.Coin),
		To:    data.To,
		Value: helpers.PipStr2Bip(data.Value),
	}
}
