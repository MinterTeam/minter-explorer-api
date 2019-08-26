package address

import (
	"github.com/MinterTeam/minter-explorer-api/balance"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"math/big"
)

type Resource struct {
	Address    string               `json:"address"`
	BalanceSum *string              `json:"balanceSum,omitempty"`
	Balances   []resource.Interface `json:"balances"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	address := model.(models.Address)

	var balanceSum *string
	if len(params) == 1 && params[0].(*big.Int) != nil {
		sum := helpers.PipStr2Bip(params[0].(*big.Int).String())
		balanceSum = &sum
	}

	return Resource{
		Address:  address.GetAddress(),
		BalanceSum: balanceSum,
		Balances: resource.TransformCollection(address.Balances, balance.Resource{}),
	}
}
