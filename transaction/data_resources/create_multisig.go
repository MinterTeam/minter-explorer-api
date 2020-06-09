package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
)

type CreateMultisig struct {
	Threshold       string   `json:"threshold"`
	Weights         []string `json:"weights"`
	Addresses       []string `json:"addresses"`
	MultisigAddress string   `json:"multisig_address"`
}

func (CreateMultisig) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*models.CreateMultisigTxData)
	model := params[0].(models.Transaction)

	return CreateMultisig{
		Threshold:       data.Threshold,
		Weights:         data.Weights,
		Addresses:       data.Addresses,
		MultisigAddress: `Mx` + model.Tags["tx.created_multisig"],
	}
}
