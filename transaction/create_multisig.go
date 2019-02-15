package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type CreateMultisigDataResource struct {
	Threshold string   `json:"threshold"`
	Weights   string   `json:"weights"`
	Addresses []string `json:"addresses"`
}

func (CreateMultisigDataResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.CreateMultisigData)

	return CreateMultisigDataResource{
		Threshold: data.Threshold,
		Weights:   data.Weights,
		Addresses: data.Addresses,
	}
}

func (resource CreateMultisigDataResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.CreateMultisigData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
