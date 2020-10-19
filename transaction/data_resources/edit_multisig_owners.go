package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditMultisigData struct {
	Threshold uint64   `json:"threshold"`
	Weights   []uint64 `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (EditMultisigData) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditMultisigData)

	return EditMultisigData{
		Threshold: data.GetThreshold(),
		Weights:   data.GetWeights(),
		Addresses: data.GetAddresses(),
	}
}
