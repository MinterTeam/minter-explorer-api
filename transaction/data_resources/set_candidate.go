package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type SetCandidate struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	var pubkey string

	if data, ok := txData.(*api_pb.SetCandidateOffData); ok {
		pubkey = data.PubKey
	}

	if data, ok := txData.(*api_pb.SetCandidateOffData); ok {
		pubkey = data.PubKey
	}

	return SetCandidate{
		PubKey: pubkey,
	}
}
