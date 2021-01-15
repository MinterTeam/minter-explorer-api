package check

import (
	"encoding/base64"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction/data_models"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	FromAddress string             `json:"address_from"`
	ToAddress   string             `json:"address_to"`
	Coin        resource.Interface `json:"coin"`
	GasCoin     resource.Interface `json:"gas_coin"`
	Nonce       string             `json:"nonce"`
	Value       string             `json:"value"`
	DueBlock    uint64             `json:"due_block"`
	RawCheck    string             `json:"raw_check"`
}

type Params struct {
	CheckData *data_models.CheckData
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	c := model.(models.Check)
	data := resourceParams[0].(Params).CheckData

	return Resource{
		FromAddress: c.FromAddress.GetAddress(),
		ToAddress:   c.ToAddress.GetAddress(),
		Coin:        new(coins.IdResource).Transform(data.Coin),
		GasCoin:     new(coins.IdResource).Transform(data.GasCoin),
		Nonce:       base64.StdEncoding.EncodeToString(data.Nonce),
		Value:       helpers.PipStr2Bip(data.Value.String()),
		DueBlock:    data.DueBlock,
		RawCheck:    `Mx` + c.Data,
	}
}
