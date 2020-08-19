package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type ChangeCoinOwner struct {
	Symbol   string `json:"symbol"`
	NewOwner string `json:"new_owner"`
}

func (ChangeCoinOwner) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.ChangeCoinOwnerData)

	return ChangeCoinOwner{
		Symbol:   data.GetSymbol(),
		NewOwner: data.GetNewOwner(),
	}
}
