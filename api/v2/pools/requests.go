package pools

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"math/big"
	"strconv"
)

type GetSwapPoolRequest struct {
	Token string `uri:"token" binding:"omitempty,required_without_all=coin0 coin1"`
	Coin0 string `uri:"coin0" binding:"omitempty,required_with=coin0"`
	Coin1 string `uri:"coin1" binding:"omitempty,required_with=coin1"`
}

type GetSwapPoolProviderRequest struct {
	Token   string `uri:"token"   binding:"omitempty,required_without_all=coin0 coin1"`
	Coin0   string `uri:"coin0"   binding:"omitempty,required_with=coin0"`
	Coin1   string `uri:"coin1"   binding:"omitempty,required_with=coin1"`
	Address string `uri:"address" binding:"minterAddress"`
}

type GetSwapPoolsRequest struct {
	Coin    *string `form:"coin"     binding:"omitempty"`
	Address *string `form:"provider" binding:"omitempty,minterAddress"`
}

type GetSwapPoolProvidersRequest struct {
	Token string `uri:"token" binding:"omitempty,required_without_all=coin0 coin1"`
	Coin0 string `uri:"coin0" binding:"omitempty,required_with=coin0"`
	Coin1 string `uri:"coin1" binding:"omitempty,required_with=coin1"`
}

type GetSwapPoolsByProviderRequest struct {
	Address string `uri:"address" binding:"required,minterAddress"`
}

type FindSwapPoolRouteRequest struct {
	Coin0 string `uri:"coin0"  binding:"required"`
	Coin1 string `uri:"coin1"  binding:"required"`
}

type GetSwapPoolOrdersRequest struct {
	Type    string `form:"type"    binding:"omitempty,oneof=buy sell"`
	Address string `form:"address" binding:"omitempty,minterAddress"`
	Status  string `form:"status"  binding:"omitempty,oneof=canceled active expired partially_filled filled"`
}

func (req *FindSwapPoolRouteRequest) GetCoins(explorer *core.Explorer) (coinFrom models.Coin, coinTo models.Coin, err error) {
	fromCoinId, toCoinId := uint64(0), uint64(0)
	if id, err := strconv.ParseUint(req.Coin0, 10, 64); err == nil {
		fromCoinId = id
	} else {
		fromCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin0)
		if err != nil {
			return coinFrom, coinTo, err
		}
	}

	if id, err := strconv.ParseUint(req.Coin1, 10, 64); err == nil {
		toCoinId = id
	} else {
		toCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin1)
		if err != nil {
			return coinFrom, coinTo, err
		}
	}

	coinFrom, err = explorer.CoinRepository.FindByID(uint(fromCoinId))
	if err != nil {
		return coinFrom, coinTo, err
	}

	coinTo, err = explorer.CoinRepository.FindByID(uint(toCoinId))
	if err != nil {
		return coinFrom, coinTo, err
	}

	return coinFrom, coinTo, err
}

type FindSwapPoolRouteRequestQuery struct {
	Amount    string `form:"amount" binding:"required,numeric"`
	TradeType string `form:"type"   binding:"required,oneof=input output"`
}

func (req *FindSwapPoolRouteRequestQuery) GetAmount() *big.Int {
	return helpers.StringToBigInt(req.Amount)
}

type GetSwapPoolTransactionsRequest struct {
	Token      string  `uri:"token"         binding:"omitempty,required_without_all=coin0 coin1"`
	Coin0      string  `uri:"coin0"         binding:"omitempty,required_with=coin0"`
	Coin1      string  `uri:"coin1"         binding:"omitempty,required_with=coin1"`
	Page       string  `form:"page"         binding:"omitempty,numeric"`
	StartBlock *string `form:"start_block"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"end_block"    binding:"omitempty,numeric"`
}

type GetCoinPossibleSwapsRequest struct {
	Coin string `uri:"coin"  binding:"required"`
}

type GetCoinPossibleSwapsRequestQuery struct {
	Depth int `form:"depth"`
}

type GetSwapPoolTradesVolumeRequestQuery struct {
	Scale *string `form:"scale" binding:"omitempty,eq=minute|eq=hour|eq=day|eq=month"`
}

type GetSwapPoolOrderRequest struct {
	OrderId uint64 `uri:"orderId" binding:"required"`
}
