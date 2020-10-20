package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditCoinOwner struct {
	Symbol   string `json:"symbol"`
	NewOwner string `json:"new_owner"`
}

func (EditCoinOwner) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditCoinOwnerData)

	return EditCoinOwner{
		Symbol:   data.GetSymbol(),
		NewOwner: data.GetNewOwner(),
	}
}
