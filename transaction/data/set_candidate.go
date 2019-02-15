package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type SetCandidateResource struct {
	PubKey string `json:"pub_key"`
}

func (SetCandidateResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.SetCandidateTxData)

	return SetCandidateResource{
		PubKey: data.PubKey,
	}
}
