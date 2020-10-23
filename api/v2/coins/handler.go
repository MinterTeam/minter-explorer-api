package coins

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const CacheBlocksCount = 1

// Get list of coins
func GetCoins(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)
	symbol, version := helpers.GetSymbolAndVersionFromStr(c.Query("symbol"))

	var data []models.Coin
	if symbol == "" {
		// fetch coins resource
		data = explorer.Cache.Get("coins", func() interface{} {
			return explorer.CoinRepository.GetCoins()
		}, CacheBlocksCount).([]models.Coin)
	} else {
		// fetch coins by symbol
		data = explorer.CoinRepository.GetLikeSymbolAndVersion(symbol, version)
	}

	// make response as empty array if no models found
	if len(data) == 0 {
		empty := make([]coins.Resource, 0)

		c.JSON(http.StatusOK, gin.H{"data": empty})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resource.TransformCollection(data, coins.Resource{})})
}

type GetCoinByIdRequest struct {
	ID uint `uri:"id" binding:"numeric"`
}

func GetCoinByID(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinByIdRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	coin, err := explorer.CoinRepository.FindByID(request.ID)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(coins.Resource).Transform(coin),
	})
}

type GetCoinBySymbolRequest struct {
	Symbol string `uri:"symbol"`
}

func GetCoinBySymbol(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetCoinBySymbolRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	symbol, version := helpers.GetSymbolAndVersionFromStr(request.Symbol)
	models := explorer.CoinRepository.GetBySymbolAndVersion(symbol, version)
	if len(models) == 0 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(coins.Resource).Transform(models[0]),
	})
}
