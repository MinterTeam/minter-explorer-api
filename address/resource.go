package address

import (
	"github.com/MinterTeam/minter-explorer-api/v2/balance"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction/data_resources"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
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
	TotalBalanceSum    *big.Int
	TotalBalanceSumUSD *big.Float
	StakeBalanceSum    *big.Int
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

type LockedTokenResource struct {
	Coin       data_resources.Coin `json:"coin"`
	DueBlock   uint64              `json:"due_block"`
	Value      string              `json:"value"`
	StartBlock uint64              `json:"start_block"`
}

func (r LockedTokenResource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	m := model.(models.Transaction)
	data := m.IData.(*api_pb.LockData)

	return LockedTokenResource{
		Coin:       new(data_resources.Coin).Transform(data.Coin),
		DueBlock:   data.DueBlock,
		Value:      helpers.PipStr2Bip(data.Value),
		StartBlock: m.BlockID,
	}
}
