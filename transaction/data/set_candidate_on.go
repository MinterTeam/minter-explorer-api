package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateOnResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateOnResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SetCandidateOnData)

	return SetCandidateOnResource{
		PubKey: data.PubKey,
	}
}
