package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type EditCandidateResource struct {
	PubKey        string `json:"pub_key"`
	RewardAddress string `json:"reward_address"`
	OwnerAddress  string `json:"owner_address"`
}

func (EditCandidateResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.EditCandidateData)

	return EditCandidateResource{
		PubKey:        data.PubKey,
		RewardAddress: data.RewardAddress,
		OwnerAddress:  data.OwnerAddress,
	}
}

func (resource EditCandidateResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.EditCandidateData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
