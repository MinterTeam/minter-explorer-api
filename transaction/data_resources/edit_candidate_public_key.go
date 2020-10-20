package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type EditCandidatePublicKey struct {
	PubKey    string `json:"pub_key"`
	NewPubKey string `json:"new_pub_key"`
}

func (EditCandidatePublicKey) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.EditCandidatePublicKeyData)

	return EditCandidatePublicKey{
		PubKey:    data.PubKey,
		NewPubKey: data.NewPubKey,
	}
}
