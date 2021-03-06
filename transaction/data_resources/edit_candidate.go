package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditCandidate struct {
	PubKey         string `json:"pub_key"`
	RewardAddress  string `json:"reward_address"`
	OwnerAddress   string `json:"owner_address"`
	ControlAddress string `json:"control_address"`
}

func (EditCandidate) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditCandidateData)

	return EditCandidate{
		PubKey:         data.PubKey,
		RewardAddress:  data.RewardAddress,
		OwnerAddress:   data.OwnerAddress,
		ControlAddress: data.ControlAddress,
	}
}
