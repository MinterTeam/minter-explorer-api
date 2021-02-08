package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
)

type Resource struct {
	Coin0          resource.Interface `json:"coin0"`
	Coin1          resource.Interface `json:"coin1"`
	Amount0        string             `json:"amount0"`
	Amount1        string             `json:"amount1"`
	Liquidity      string             `json:"liquidity"`
	LiquidityInBip string             `json:"liquidity_bip"`
}

type Params struct {
	LiquidityInBip *big.Int
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	pool := model.(models.LiquidityPool)

	return Resource{
		Coin0:          new(coins.IdResource).Transform(*pool.FirstCoin),
		Coin1:          new(coins.IdResource).Transform(*pool.SecondCoin),
		Amount0:        helpers.PipStr2Bip(pool.FirstCoinVolume),
		Amount1:        helpers.PipStr2Bip(pool.SecondCoinVolume),
		Liquidity:      helpers.PipStr2Bip(pool.Liquidity),
		LiquidityInBip: helpers.PipStr2Bip(resourceParams[0].(Params).LiquidityInBip.String()),
	}
}

type ProviderResource struct {
	Coin0          resource.Interface `json:"coin0"`
	Coin1          resource.Interface `json:"coin1"`
	Amount0        string             `json:"amount0"`
	Amount1        string             `json:"amount1"`
	Liquidity      string             `json:"liquidity"`
	LiquidityInBip string             `json:"liquidity_bip"`
}

func (r ProviderResource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	provider := model.(models.AddressLiquidityPool)

	part := new(big.Float).Quo(
		new(big.Float).SetInt(helpers.StringToBigInt(provider.Liquidity)),
		new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.Liquidity)),
	)

	amount0, _ := new(big.Float).Mul(part, new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.FirstCoinVolume))).Int(nil)
	amount1, _ := new(big.Float).Mul(part, new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.SecondCoinVolume))).Int(nil)

	liquidityInBip, _ := new(big.Float).Mul(part, new(big.Float).SetInt(resourceParams[0].(Params).LiquidityInBip)).Int(nil)

	return ProviderResource{
		Coin0:          new(coins.IdResource).Transform(*provider.LiquidityPool.FirstCoin),
		Coin1:          new(coins.IdResource).Transform(*provider.LiquidityPool.SecondCoin),
		Amount0:        helpers.PipStr2Bip(amount0.String()),
		Amount1:        helpers.PipStr2Bip(amount1.String()),
		Liquidity:      helpers.PipStr2Bip(provider.Liquidity),
		LiquidityInBip: helpers.PipStr2Bip(liquidityInBip.String()),
	}
}
