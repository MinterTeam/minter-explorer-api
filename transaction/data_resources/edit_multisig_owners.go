package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditMultisigOwners struct {
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (EditMultisigOwners) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditMultisigOwnersData)

	return EditMultisigOwners{
		Threshold: data.GetThreshold(),
		Weights:   data.GetWeights(),
		Addresses: data.GetAddresses(),
	}
}
