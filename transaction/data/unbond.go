package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type UnbondResource struct {
	PubKey string `json:"pub_key"`
	Coin   string `json:"coin"`
	Value  string `json:"value"`
}

func (UnbondResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.UnbondData)

	return UnbondResource{
		PubKey: data.PubKey,
		Coin:   data.Coin,
		Value:  helpers.PipStr2Bip(data.Value),
	}
}

func (resource UnbondResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.UnbondData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
