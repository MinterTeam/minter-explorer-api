package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type RedeemCheckResource struct {
	RawCheck string `json:"raw_check"`
	Proof    string `json:"proof"`
}

func (RedeemCheckResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.RedeemCheckData)

	return RedeemCheckResource{
		RawCheck: data.RawCheck,
		Proof:    data.Proof,
	}
}

func (resource RedeemCheckResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.RedeemCheckData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
