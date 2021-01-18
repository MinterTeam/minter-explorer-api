package pools

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetSwapPoolRequest struct {
	Coin0 uint64 `uri:"coin0" binding:"numeric"`
	Coin1 uint64 `uri:"coin1" binding:"numeric"`
}

func GetSwapPool(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(req.Coin0, req.Coin1)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pool not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.Resource).Transform(p),
	})
}

type GetSwapPoolProviderRequest struct {
	Coin0   uint64 `uri:"coin0" binding:"numeric"`
	Coin1   uint64 `uri:"coin1" binding:"numeric"`
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

	p, err := explorer.PoolRepository.FindProvider(req.Coin0, req.Coin1, helpers.RemoveMinterPrefix(req.Address))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.ProviderResource).Transform(p),
	})
}

type GetSwapPoolsRequest struct {
	Coin    *uint64 `form:"coin_id" binding:"omitempty,numeric"`
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
		CoinId:          req.Coin,
		ProviderAddress: req.Address,
	}, &pagination)

	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(pools, pool.Resource{}, pagination))
}
