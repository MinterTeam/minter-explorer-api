package transaction

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type DeclareCandidacyDataResource struct {
	Address    string `json:"address"`
	PubKey     string `json:"pub_key"`
	Commission string `json:"commission"`
	Coin       string `json:"coin"`
	Stake      string `json:"stake"`
}

func (DeclareCandidacyDataResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.DeclareCandidacyData)

	return DeclareCandidacyDataResource{
		Address:    data.Address,
		PubKey:     data.PubKey,
		Commission: data.Commission,
		Coin:       data.Coin,
		Stake:      helpers.PipStr2Bip(data.Stake),
	}
}

func (resource DeclareCandidacyDataResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.DeclareCandidacyData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
