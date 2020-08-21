package address

import (
	"github.com/MinterTeam/minter-explorer-api/balance"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
)

type Resource struct {
	Address            string               `json:"address"`
	Balances           []resource.Interface `json:"balances"`
	TotalBalanceSum    *string              `json:"total_balance_sum,omitempty"`
	TotalBalanceSumUSD *string              `json:"total_balance_sum_usd,omitempty"`
	StakeBalanceSum    *string              `json:"stake_balance_sum,omitempty"`
	StakeBalanceSumUSD *string              `json:"stake_balance_sum_usd,omitempty"`
}

type Params struct {
	TotalBalanceSum    *big.Float
	TotalBalanceSumUSD *big.Float
	StakeBalanceSum    *big.Float
	StakeBalanceSumUSD *big.Float
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	address := model.(models.Address)

	r.Address = address.GetAddress()
	r.Balances = resource.TransformCollection(address.Balances, balance.Resource{})

	if len(resourceParams) > 0 {
		if params, ok := resourceParams[0].(Params); ok {
			r.transformParams(params)
		}
	}

	return r
}

// prepare total address balance
func (r *Resource) transformParams(params Params) {
	sum := helpers.PipStr2Bip(params.TotalBalanceSum.String())
	usd := helpers.PipStr2Bip(params.TotalBalanceSumUSD.String())

	stakeSum := helpers.PipStr2Bip(params.StakeBalanceSum.String())
	stakeSumUSD := helpers.PipStr2Bip(params.StakeBalanceSumUSD.String())

	r.TotalBalanceSum, r.TotalBalanceSumUSD = &sum, &usd
	r.StakeBalanceSum, r.StakeBalanceSumUSD = &stakeSum, &stakeSumUSD
}
