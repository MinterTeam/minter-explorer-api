package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type LockData struct {
	DueBlock uint64 `json:"due_block"`
	Coin     Coin   `json:"coin"`
	Value    string `json:"value"`
}

func (LockData) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.LockData)

	return LockData{
		DueBlock: data.DueBlock,
		Coin:     new(Coin).Transform(data.Coin),
		Value:    helpers.PipStr2Bip(data.Value),
	}
}
