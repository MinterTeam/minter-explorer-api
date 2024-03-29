package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
	"strconv"
	"time"
)

type Resource struct {
	Coin0             resource.Interface `json:"coin0"`
	Coin1             resource.Interface `json:"coin1"`
	Token             resource.Interface `json:"token"`
	Amount0           string             `json:"amount0"`
	Amount1           string             `json:"amount1"`
	Liquidity         string             `json:"liquidity"`
	LiquidityInBip    string             `json:"liquidity_bip"`
	TradeVolumeBip1d  string             `json:"trade_volume_bip_1d"`
	TradeVolumeBip30d string             `json:"trade_volume_bip_30d"`
	UpdatedAtBlock    uint64             `json:"updated_at_block"`
}

type Params struct {
	TradeVolume1d  *big.Float
	TradeVolume30d *big.Float
	FirstCoin      string
	SecondCoin     string
}

func (r Resource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	pool := model.(models.LiquidityPool)
	params := resourceParams[0].(Params)

	inverseOrder := false
	if len(params.FirstCoin) != 0 && len(params.SecondCoin) != 0 {
		if id, err := strconv.Atoi(params.FirstCoin); err == nil {
			inverseOrder = uint(id) != pool.FirstCoin.ID && params.FirstCoin != pool.FirstCoin.GetSymbol()
		} else if params.FirstCoin != pool.FirstCoin.GetSymbol() {
			inverseOrder = true
		}
	}

	if !inverseOrder {
		return Resource{
			Coin0:             new(coins.IdResource).Transform(*pool.FirstCoin),
			Coin1:             new(coins.IdResource).Transform(*pool.SecondCoin),
			Token:             new(coins.IdResource).Transform(*pool.Token),
			Amount0:           helpers.PipStr2Bip(pool.FirstCoinVolume),
			Amount1:           helpers.PipStr2Bip(pool.SecondCoinVolume),
			Liquidity:         helpers.PipStr2Bip(pool.Liquidity),
			LiquidityInBip:    helpers.Pip2BipStr(helpers.StringToBigInt(pool.LiquidityBip)),
			TradeVolumeBip1d:  helpers.Bip2Str(params.TradeVolume1d),
			TradeVolumeBip30d: helpers.Bip2Str(params.TradeVolume30d),
			UpdatedAtBlock:    pool.UpdatedAtBlockId,
		}
	}

	return Resource{
		Coin0:             new(coins.IdResource).Transform(*pool.SecondCoin),
		Coin1:             new(coins.IdResource).Transform(*pool.FirstCoin),
		Token:             new(coins.IdResource).Transform(*pool.Token),
		Amount0:           helpers.PipStr2Bip(pool.SecondCoinVolume),
		Amount1:           helpers.PipStr2Bip(pool.FirstCoinVolume),
		Liquidity:         helpers.PipStr2Bip(pool.Liquidity),
		LiquidityInBip:    helpers.Pip2BipStr(helpers.StringToBigInt(pool.LiquidityBip)),
		TradeVolumeBip1d:  helpers.Bip2Str(params.TradeVolume1d),
		TradeVolumeBip30d: helpers.Bip2Str(params.TradeVolume30d),
		UpdatedAtBlock:    pool.UpdatedAtBlockId,
	}
}

type ProviderResource struct {
	Address        string             `json:"address"`
	Coin0          resource.Interface `json:"coin0"`
	Coin1          resource.Interface `json:"coin1"`
	Token          resource.Interface `json:"token"`
	Amount0        string             `json:"amount0"`
	Amount1        string             `json:"amount1"`
	Liquidity      string             `json:"liquidity"`
	LiquidityShare float64            `json:"liquidity_share"`
	LiquidityInBip string             `json:"liquidity_bip"`
}

func (r ProviderResource) Transform(model resource.ItemInterface, resourceParams ...resource.ParamInterface) resource.Interface {
	provider := model.(models.AddressLiquidityPool)
	params := resourceParams[0].(Params)

	part := new(big.Float).Quo(
		new(big.Float).SetInt(helpers.StringToBigInt(provider.Liquidity)),
		new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.Liquidity)),
	)

	amount0, _ := new(big.Float).Mul(part, new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.FirstCoinVolume))).Int(nil)
	amount1, _ := new(big.Float).Mul(part, new(big.Float).SetInt(helpers.StringToBigInt(provider.LiquidityPool.SecondCoinVolume))).Int(nil)

	liquidityShare, _ := part.Float64()
	liquidityInBip := new(big.Float).Mul(part, helpers.Pip2Bip(helpers.StringToBigInt(provider.LiquidityPool.LiquidityBip)))

	inverseOrder := false

	if len(params.FirstCoin) != 0 && len(params.SecondCoin) != 0 {
		if id, err := strconv.Atoi(params.FirstCoin); err == nil {
			inverseOrder = uint(id) != provider.LiquidityPool.FirstCoin.ID
		} else if params.FirstCoin != provider.LiquidityPool.FirstCoin.GetSymbol() {
			inverseOrder = true
		}
	}

	if !inverseOrder {
		return ProviderResource{
			Address:        provider.Address.GetAddress(),
			Coin0:          new(coins.IdResource).Transform(*provider.LiquidityPool.FirstCoin),
			Coin1:          new(coins.IdResource).Transform(*provider.LiquidityPool.SecondCoin),
			Token:          new(coins.IdResource).Transform(*provider.LiquidityPool.Token),
			Amount0:        helpers.PipStr2Bip(amount0.String()),
			Amount1:        helpers.PipStr2Bip(amount1.String()),
			Liquidity:      helpers.PipStr2Bip(provider.Liquidity),
			LiquidityInBip: helpers.Bip2Str(liquidityInBip),
			LiquidityShare: liquidityShare * 100,
		}
	}

	return ProviderResource{
		Address:        provider.Address.GetAddress(),
		Coin0:          new(coins.IdResource).Transform(*provider.LiquidityPool.SecondCoin),
		Coin1:          new(coins.IdResource).Transform(*provider.LiquidityPool.FirstCoin),
		Token:          new(coins.IdResource).Transform(*provider.LiquidityPool.Token),
		Amount0:        helpers.PipStr2Bip(amount1.String()),
		Amount1:        helpers.PipStr2Bip(amount0.String()),
		Liquidity:      helpers.PipStr2Bip(provider.Liquidity),
		LiquidityInBip: helpers.Bip2Str(liquidityInBip),
		LiquidityShare: liquidityShare * 100,
	}
}

type BancorResource struct {
	SwapType  string `json:"swap_type"`
	AmountIn  string `json:"amount_in"`
	AmountOut string `json:"amount_out"`
}

func (r BancorResource) Transform(value *big.Int, bancor *big.Int, tradeType TradeType) BancorResource {
	if tradeType == TradeTypeExactInput {
		return BancorResource{
			SwapType:  "bancor",
			AmountIn:  helpers.Pip2BipStr(value),
			AmountOut: helpers.Pip2BipStr(bancor),
		}
	}

	return BancorResource{
		SwapType:  "bancor",
		AmountIn:  helpers.Pip2BipStr(bancor),
		AmountOut: helpers.Pip2BipStr(value),
	}
}

type TradesVolumeResource struct {
	Date             string `json:"date"`
	FirstCoinVolume  string `json:"first_coin_volume"`
	SecondCoinVolume string `json:"second_coin_volume"`
	BipVolume        string `json:"bip_volume"`
}

func (r TradesVolumeResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	tv := model.(TradeVolume)

	return TradesVolumeResource{
		Date:             tv.Date.Format(time.RFC3339),
		FirstCoinVolume:  helpers.PipStr2Bip(tv.FirstCoinVolume),
		SecondCoinVolume: helpers.PipStr2Bip(tv.SecondCoinVolume),
		BipVolume:        helpers.Bip2Str(tv.BipVolume),
	}
}
