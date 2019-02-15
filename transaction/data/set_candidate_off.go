package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateOffResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateOffResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SetCandidateOffData)

	return SetCandidateOffResource{
		PubKey: data.PubKey,
	}
}
