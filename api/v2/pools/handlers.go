package pools

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetSwapPoolRequest struct {
	Coin0 string `uri:"coin0"`
	Coin1 string `uri:"coin1"`
}

func GetSwapPool(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1})
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
	Coin0   string `uri:"coin0" binding:""`
	Coin1   string `uri:"coin1" binding:""`
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

	p, err := explorer.PoolRepository.FindProvider(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1}, helpers.RemoveMinterPrefix(req.Address))
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

	if len(pools) == 0 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pools not found.", c)
		return
	}

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(model.(models.LiquidityPool))
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(pools, pool.Resource{}, pagination, resourceCallback))
}

type GetSwapPoolProvidersRequest struct {
	Coin0 string `uri:"coin0" binding:""`
	Coin1 string `uri:"coin1" binding:""`
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
	providers, err := explorer.PoolRepository.GetProviders(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1}, &pagination)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
		return
	}

	if len(providers) == 0 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Providers not found.", c)
		return
	}

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		bipValue := explorer.PoolService.GetPoolLiquidityInBip(*model.(models.AddressLiquidityPool).LiquidityPool)
		return resource.ParamsInterface{pool.Params{LiquidityInBip: bipValue}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(providers, pool.ProviderResource{}, pagination, resourceCallback))
}
