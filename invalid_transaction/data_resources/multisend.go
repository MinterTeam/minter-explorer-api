package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type Multisend struct {
	List []Send `json:"list"`
}

func (Multisend) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.MultiSendData)

	list := make([]Send, len(data.List))
	for key, item := range data.List {
		list[key] = Send{}.Transform(item).(Send)
	}

	return Multisend{list}
}
