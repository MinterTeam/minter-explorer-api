package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateOffDataResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateOffDataResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SetCandidateOffData)

	return SetCandidateOffDataResource{
		PubKey: data.PubKey,
	}
}

func (resource SetCandidateOffDataResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SetCandidateOffData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
