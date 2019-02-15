package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type MultisendDataResource struct {
	List []SendDataResource `json:"list"`
}

func (MultisendDataResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.MultisendData)

	var list []SendDataResource
	for _, item := range data.List {
		list = append(list, SendDataResource{}.Transform(item).(SendDataResource))
	}

	return MultisendDataResource{
		List: list,
	}
}

func (resource MultisendDataResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.MultisendData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
