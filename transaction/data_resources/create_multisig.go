package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type CreateMultisig struct {
	Threshold string   `json:"threshold"`
	Weights   []string `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (CreateMultisig) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.CreateMultisigTxData)

	return CreateMultisig{
		Threshold: data.Threshold,
		Weights:   data.Weights,
		Addresses: data.Addresses,
	}
}
