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

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.Resource).Transform(p),
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

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.ProviderResource).Transform(p),
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

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(pools, pool.Resource{}, pagination))
}
