package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateOnResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateOnResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.SetCandidateOnData)

	return SetCandidateOnResource{
		PubKey: data.PubKey,
	}
}

func (resource SetCandidateOnResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.SetCandidateOnData
	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
