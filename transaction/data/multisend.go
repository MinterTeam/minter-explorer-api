package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type MultisendResource struct {
	List []SendResource `json:"list"`
}

func (MultisendResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.MultisendData)

	var list []SendResource
	for _, item := range data.List {
		list = append(list, SendResource{}.Transform(item).(SendResource))
	}

	return MultisendResource{
		List: list,
	}
}

func (resource MultisendResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.MultisendData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
