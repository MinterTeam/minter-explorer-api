package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type MoveStake struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Coin  Coin   `json:"coin"`
	Stake string `json:"stake"`
}

func (MoveStake) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.MoveStakeData)

	return MoveStake{
		From:  data.From,
		To:    data.To,
		Coin:  new(Coin).Transform(data.Coin),
		Stake: helpers.PipStr2Bip(data.Stake),
	}
}
