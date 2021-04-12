package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditCandidateCommission struct {
	PubKey     string `json:"pub_key"`
	Commission uint64 `json:"commission"`
}

func (EditCandidateCommission) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditCandidateCommission)

	return EditCandidateCommission{
		PubKey:     data.PubKey,
		Commission: data.Commission,
	}
}
