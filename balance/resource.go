package balance

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	balance := model.(models.Balance)

	return Resource{
		Coin:   balance.Coin.Symbol,
		Amount: helpers.PipStr2Bip(balance.Value),
	}
}
