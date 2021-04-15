package pools

import (
	"net/http"
	"strconv"

	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
)

type CachePoolsList struct {
	Pools      []models.LiquidityPool
	Pagination tools.Pagination
}

const (
	CachePoolCoinsBlockCount = 1
	CachePoolsListBlockCount = 1
)

func GetSwapPool(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token})
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pool not found.", c)
		return
	}

	bipValue := explorer.PoolService.GetPoolLiquidityInBip(p)

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.Resource).Transform(p, pool.Params{LiquidityInBip: bipValue, FirstCoin: req.Coin0, SecondCoin: req.Coin1}),
	})
}

func GetSwapPoolProvider(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolProviderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindProvider(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token}, helpers.RemoveMinterPrefix(req.Address))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
		return
	}

	bipValue := explorer.PoolService.GetPoolLiquidityInBip(*p.LiquidityPool)

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.ProviderResource).Transform(p, pool.Params{LiquidityInBip: bipValue, FirstCoin: req.Coin0, SecondCoin: req.Coin1}),
	})
}

func GetSwapPools(c *gin.Context) {
	var err error
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolsRequest
	if err = c.ShouldBindQuery(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// cache pools
	var pools []models.LiquidityPool

	// prepare fetching pools
	pagination := tools.NewPagination(c.Request)
	fetchPools := func() []models.LiquidityPool {
		pools, _ = explorer.PoolRepository.GetPools(pool.SelectPoolsFilter{
			Coin:            req.Coin,
			ProviderAddress: req.Address,
		}, &pagination)

		return pools
	}

	if len(c.Request.URL.Query()) == 0 {
		cached := explorer.Cache.Get("pools", func() interface{} {
			return CachePoolsList{fetchPools(), pagination}
		}, CachePoolsListBlockCount).(CachePoolsList)

		pools = cached.Pools
		pagination = cached.Pagination
	} else {
		pools = fetchPools()
	}

	helpers.CheckErr(err)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(model.(models.LiquidityPool))
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(pools, pool.Resource{}, pagination, resourceCallback))
}

func GetSwapPoolProviders(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolProvidersRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	providers, err := explorer.PoolRepository.GetProviders(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token}, &pagination)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
		return
	}

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(*model.(models.AddressLiquidityPool).LiquidityPool)
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue, FirstCoin: req.Coin0, SecondCoin: req.Coin1}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(providers, pool.ProviderResource{}, pagination, resourceCallback))
}

func GetSwapPoolsByProvider(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolsByProviderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	pools, err := explorer.PoolRepository.GetPoolsByProvider(helpers.RemoveMinterPrefix(req.Address), &pagination)
	helpers.CheckErr(err)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(*model.(models.AddressLiquidityPool).LiquidityPool)
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(pools, pool.ProviderResource{}, pagination, resourceCallback))
}

func FindSwapPoolRoute(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	if req.Coin0 == req.Coin1 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	var reqQuery FindSwapPoolRouteRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	fromCoinId, toCoinId := uint64(0), uint64(0)
	if id, err := strconv.ParseUint(req.Coin0, 10, 64); err == nil {
		fromCoinId = id
	} else {
		fromCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin0)
		if err != nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
			return
		}
	}

	if id, err := strconv.ParseUint(req.Coin1, 10, 64); err == nil {
		toCoinId = id
	} else {
		toCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin1)
		if err != nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
			return
		}
	}

	// define trade type
	tradeType := swap.TradeTypeExactInput
	if reqQuery.TradeType == "output" {
		tradeType = swap.TradeTypeExactOutput
	}

	trade, err := explorer.PoolService.FindSwapRoutePath(fromCoinId, toCoinId, tradeType, helpers.StringToBigInt(reqQuery.Amount))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		coin, err := explorer.CoinRepository.FindByID(uint(t.CoinID))
		helpers.CheckErr(err)
		path[i] = coin
	}

	c.JSON(http.StatusOK, new(pool.RouteResource).Transform(path, trade))
}

func GetSwapPoolTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolTransactionsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.PoolsFilter{
		Coin0:      req.Coin0,
		Coin1:      req.Coin1,
		Token:      req.Token,
		StartBlock: req.StartBlock,
		EndBlock:   req.EndBlock,
	}, &pagination)

	txs, err := explorer.TransactionService.PrepareTransactionsModel(txs)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

func GetCoinsList(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	poolCoins := explorer.Cache.Get("pools_coins_list", func() interface{} {
		poolCoins, _ := explorer.PoolRepository.GetPoolsCoins()
		return poolCoins
	}, CachePoolCoinsBlockCount).([]models.Coin)

	c.JSON(http.StatusOK, resource.TransformCollection(poolCoins, coins.IdResource{}))
}

func GetCoinPossibleSwaps(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetCoinPossibleSwapsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var reqQuery GetCoinPossibleSwapsRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	poolCoins := explorer.Cache.Get("pools_coins_list", func() interface{} {
		poolCoins, _ := explorer.PoolRepository.GetPoolsCoins()
		return poolCoins
	}, CachePoolCoinsBlockCount).([]models.Coin)

	// find coin from by id or symbol
	var fromCoin *models.Coin
	if id, err := strconv.ParseUint(req.Coin, 10, 64); err == nil {
		for _, pc := range poolCoins {
			if uint64(pc.ID) == id {
				fromCoin = &pc
			}
		}
	} else {
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(req.Coin)
		for _, pc := range poolCoins {
			if pc.Symbol == symbol && uint64(pc.Version) == version {
				fromCoin = &models.Coin{
					ID:     pc.ID,
					Symbol: pc.Symbol,
				}
			}
		}
	}

	// send empty response if coin not exists in pools
	if fromCoin == nil {
		c.JSON(http.StatusOK, resource.TransformCollection([]models.Coin{}, coins.IdResource{}))
		return
	}

	liquidityPools, err := explorer.PoolRepository.GetAll()
	helpers.CheckErr(err)

	var swapCoins []models.Coin
	for _, pc := range poolCoins {
		if pc.ID == fromCoin.ID {
			continue
		}

		if _, err := explorer.PoolService.FindSwapRoutePathsByGraph(liquidityPools, uint64(fromCoin.ID), uint64(pc.ID), reqQuery.Depth); err == nil {
			swapCoins = append(swapCoins, pc)
		}
	}

	c.JSON(http.StatusOK, resource.TransformCollection(swapCoins, coins.IdResource{}))
}

func EstimateSwap(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var reqQuery FindSwapPoolRouteRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	coinFrom, coinTo, err := req.GetCoins(explorer)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	// define trade type
	tradeType := swap.TradeTypeExactInput
	if reqQuery.TradeType == "output" {
		tradeType = swap.TradeTypeExactOutput
	}

	bancorAmount, bancorErr := explorer.SwapService.EstimateByBancor(coinFrom, coinTo, reqQuery.GetAmount(), tradeType)
	trade, poolErr := explorer.PoolService.FindSwapRoutePath(uint64(coinFrom.ID), uint64(coinTo.ID), tradeType, reqQuery.GetAmount())

	if poolErr != nil && bancorErr != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	if bancorErr == nil && poolErr != nil {
		c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
		return
	}

	if bancorErr == nil && poolErr == nil {
		if tradeType == swap.TradeTypeExactInput && bancorAmount.Cmp(trade.OutputAmount.GetAmount()) >= 1 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}

		if tradeType == swap.TradeTypeExactOutput && bancorAmount.Cmp(trade.InputAmount.GetAmount()) <= 0 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}
	}

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		path[i], _ = explorer.CoinRepository.FindByID(uint(t.CoinID))
	}

	c.JSON(http.StatusOK, new(pool.RouteResource).Transform(path, trade))
}

func GetSwapPoolTradesVolume(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var query GetSwapPoolTradesVolumeRequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token})
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pool not found.", c)
		return
	}

	tradesVolume, err := explorer.PoolService.GetTradesVolume(p, query.Scale)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(tradesVolume, pool.TradesVolumeResource{}),
	})
}
