package data

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type CreateCoinResource struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	InitialAmount        string `json:"initial_amount"`
	InitialReserve       string `json:"initial_reserve"`
	ConstantReserveRatio string `json:"constant_reserve_ratio"`
}

func (CreateCoinResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(models.CreateCoinData)

	return CreateCoinResource{
		Name:                 data.Name,
		Symbol:               data.Symbol,
		InitialAmount:        helpers.PipStr2Bip(data.InitialAmount),
		InitialReserve:       helpers.PipStr2Bip(data.InitialReserve),
		ConstantReserveRatio: data.ConstantReserveRatio,
	}
}

func (resource CreateCoinResource) TransformFromJsonRaw(raw json.RawMessage) resource.Interface {
	var data models.CreateCoinData

	err := json.Unmarshal(raw, &data)
	helpers.CheckErr(err)

	return resource.Transform(data)
}
