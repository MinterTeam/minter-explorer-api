package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type CreateMultisigResource struct {
	Threshold string   `json:"threshold"`
	Weights   string   `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (CreateMultisigResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.CreateMultisigTxData)

	return CreateMultisigResource{
		Threshold: data.Threshold,
		Weights:   data.Weights,
		Addresses: data.Addresses,
	}
}
