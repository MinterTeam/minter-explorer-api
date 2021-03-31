package balance

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/services"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Coin      resource.Interface `json:"coin"`
	Amount    string             `json:"amount"`
	BipAmount string             `json:"bip_amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)
	bipAmount := services.Swap.EstimateInBip(balance.Coin, helpers.StringToBigInt(balance.Value))

	return Resource{
		Coin:      new(coins.IdResource).Transform(*balance.Coin, coins.Params{IsTypeRequired: true}),
		Amount:    helpers.PipStr2Bip(balance.Value),
		BipAmount: helpers.Pip2BipStr(bipAmount),
	}
}
