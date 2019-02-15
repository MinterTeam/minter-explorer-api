package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateOffResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateOffResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SetCandidateOffData)

	return SetCandidateOffResource{
		PubKey: data.PubKey,
	}
}

func (resource SetCandidateOffResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SetCandidateOffData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
