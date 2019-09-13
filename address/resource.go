package address

import (
	"github.com/MinterTeam/minter-explorer-api/balance"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"math/big"
)

type Resource struct {
	Address              string               `json:"address"`
	BalanceSumInBaseCoin *string              `json:"balanceSumInBaseCoin,omitempty"`
	BalanceSumInUSD      *string              `json:"balanceSumInUSD,omitempty"`
	Balances             []resource.Interface `json:"balances"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	address := model.(models.Address)
	balanceSumInBaseCoin, balanceSumInUSD := getBalanceSumParams(params...)

	return Resource{
		Address:              address.GetAddress(),
		BalanceSumInBaseCoin: balanceSumInBaseCoin,
		BalanceSumInUSD:      balanceSumInUSD,
		Balances:             resource.TransformCollection(address.Balances, balance.Resource{}),
	}
}

func getBalanceSumParams(params ...resource.ParamInterface) (*string, *string) {
	if len(params) != 2 {
		return nil, nil
	}

	var sumInUSD, sumInBaseCoin string

	if sum, ok := params[0].(*big.Float); ok && sum != nil {
		sumInBaseCoin = helpers.PipStr2Bip(sum.String())
	} else {
		return nil, nil
	}

	if sum, ok := params[1].(*big.Float); ok && sum != nil {
		sumInUSD = helpers.PipStr2Bip(sum.String())
	} else {
		return nil, nil
	}

	return &sumInBaseCoin, &sumInUSD
}
