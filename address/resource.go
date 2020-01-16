package address

import (
	"github.com/MinterTeam/minter-explorer-api/balance"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"math/big"
)

type Resource struct {
	Address                string               `json:"address"`
	Balances               []resource.Interface `json:"balances"`
	TotalBalanceSum        *string              `json:"total_balance_sum,omitempty"`
	TotalBalanceSumUSD     *string              `json:"total_balance_sum_usd,omitempty"`
}

type Params struct {
	TotalBalanceSum    *big.Float
	TotalBalanceSumUSD *big.Float
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	address := model.(models.Address)
	result := Resource{
		Address:  address.GetAddress(),
		Balances: resource.TransformCollection(address.Balances, balance.Resource{}),
	}

	if len(resourceParams) > 0 {
		if params, ok := resourceParams[0].(Params); ok {
			result.TotalBalanceSum, result.TotalBalanceSumUSD = r.getTotalBalanceParams(params)
		}
	}

	return result
}

// prepare total address balance
func (r Resource) getTotalBalanceParams(params Params) (*string, *string) {
	sum := helpers.PipStr2Bip(params.TotalBalanceSum.String())
	usd := helpers.PipStr2Bip(params.TotalBalanceSumUSD.String())

	return &sum, &usd
}
