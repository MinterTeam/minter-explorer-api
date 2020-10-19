package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type CreateMultisig struct {
	Threshold       uint64   `json:"threshold"`
	Weights         []uint64 `json:"weights"`
	Addresses       []string `json:"addresses"`
	MultisigAddress string   `json:"multisig_address"`
}

func (CreateMultisig) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.CreateMultisigData)
	model := params[0].(models.Transaction)

	return CreateMultisig{
		Threshold:       data.Threshold,
		Weights:         data.Weights,
		Addresses:       data.Addresses,
		MultisigAddress: `Mx` + model.Tags["tx.created_multisig"],
	}
}
