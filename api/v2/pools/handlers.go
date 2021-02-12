package pools

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetSwapPoolRequest struct {
	Token string `uri:"token" validate:"required_without_all=coin0 coin1"`
	Coin0 string `uri:"coin0" validate:"required_with=coin0"`
	Coin1 string `uri:"coin1" validate:"required_with=coin1"`
}

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
		"data": new(pool.Resource).Transform(p, pool.Params{LiquidityInBip: bipValue}),
	})
}

type GetSwapPoolProviderRequest struct {
	Token   string `uri:"token" validate:"required_without_all=coin0 coin1"`
	Coin0   string `uri:"coin0" validate:"required_with=coin0"`
	Coin1   string `uri:"coin1" validate:"required_with=coin1"`
	Address string `uri:"address" binding:"minterAddress"`
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
		"data": new(pool.ProviderResource).Transform(p, pool.Params{LiquidityInBip: bipValue}),
	})
}

type GetSwapPoolsRequest struct {
	Coin    *string `form:"coin"     binding:"omitempty"`
	Address *string `form:"provider" binding:"omitempty,minterAddress"`
}

func GetSwapPools(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	pools, err := explorer.PoolRepository.GetPools(pool.SelectPoolsFilter{
		Coin:            req.Coin,
		ProviderAddress: req.Address,
	}, &pagination)

	helpers.CheckErr(err)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(model.(models.LiquidityPool))
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(pools, pool.Resource{}, pagination, resourceCallback))
}

type GetSwapPoolProvidersRequest struct {
	Token string `uri:"token" validate:"required_without_all=coin0 coin1"`
	Coin0 string `uri:"coin0" validate:"required_with=coin0"`
	Coin1 string `uri:"coin1" validate:"required_with=coin1"`
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
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(providers, pool.ProviderResource{}, pagination, resourceCallback))
}

type GetSwapPoolsByProviderRequest struct {
	Address string `uri:"address" binding:"required,minterAddress"`
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

type FindSwapPoolRouteRequest struct {
	Coin0 string `uri:"coin0" validate:"required_with=coin0"`
	Coin1 string `uri:"coin1" validate:"required_with=coin1"`
}

func FindSwapPoolRoute(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pools, err := explorer.PoolRepository.FindRoutePath(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1})
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(model.(models.LiquidityPool))
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformCollectionWithCallback(pools, pool.Resource{}, resourceCallback))
}

type GetSwapPoolTransactionsRequest struct {
	Token      string  `uri:"token" validate:"required_without_all=coin0 coin1"`
	Coin0      string  `uri:"coin0" validate:"required_with=coin0"`
	Coin1      string  `uri:"coin1" validate:"required_with=coin1"`
	Page       string  `form:"page"         binding:"omitempty,numeric"`
	StartBlock *string `form:"start_block"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"end_block"    binding:"omitempty,numeric"`
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
