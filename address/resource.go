package address

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Resource struct {
	Address  string            `json:"address"`
	Balances []BalanceResource `json:"balances"`
}

type BalanceResource struct {
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

func (Resource) Transform(model resource.ItemInterface) resource.Interface {
	address := model.(models.Address)

	var balances []BalanceResource
	for _, balance := range address.Balances {
		balances = append(balances, BalanceResource{
			Coin:   balance.Coin.Symbol,
			Amount: helpers.PipStr2Bip(balance.Value.String()),
		})
	}

	// make empty array if no models
	if len(balances) == 0 {
		balances = make([]BalanceResource, 0)
	}

	return Resource{
		Address:  address.GetAddress(),
		Balances: balances,
	}
}
