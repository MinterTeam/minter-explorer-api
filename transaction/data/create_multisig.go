package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type CreateMultisigResource struct {
	Threshold string   `json:"threshold"`
	Weights   string   `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (CreateMultisigResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.CreateMultisigData)

	return CreateMultisigResource{
		Threshold: data.Threshold,
		Weights:   data.Weights,
		Addresses: data.Addresses,
	}
}

func (resource CreateMultisigResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.CreateMultisigData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
