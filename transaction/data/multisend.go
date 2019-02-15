package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type MultisendResource struct {
	List []SendResource `json:"list"`
}

func (MultisendResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.MultisendData)

	var list []SendResource
	for _, item := range data.List {
		list = append(list, SendResource{}.Transform(&item).(SendResource))
	}

	return MultisendResource{
		List: list,
	}
}
