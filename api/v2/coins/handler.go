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

	// fetch coins resource
	if symbol == "" {
		data := explorer.Cache.Get("coins", func() interface{} {
			return resource.TransformCollectionWithCallback(
				explorer.CoinRepository.GetCoins(), coins.Resource{}, extendResourcesWithTradingVolumesCallback(explorer),
			)
		}, CacheBlocksCount).([]resource.Interface)

		c.JSON(http.StatusOK, gin.H{"data": data})
		return
	}

	// fetch coins by symbol
	data := explorer.CoinRepository.GetLikeSymbolAndVersion(symbol, version)

	// make response as empty array if no models found
	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": make([]coins.Resource, 0)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resource.TransformCollectionWithCallback(
		data, coins.Resource{}, extendResourcesWithTradingVolumesCallback(explorer),
	)})
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

	coin, err := explorer.CoinRepository.FindByIdWithOwner(request.ID)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(coins.Resource).Transform(coin, extendResourceWithTradingVolumesParams(explorer, coin)),
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

	symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(request.Symbol)
	coinModels := explorer.CoinRepository.GetBySymbolAndVersion(symbol, &version)
	if len(coinModels) == 0 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(coins.Resource).Transform(coinModels[0], extendResourceWithTradingVolumesParams(explorer, coinModels[0])),
	})
}

func GetOracleVerifiedCoins(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)
	verified := explorer.CoinService.GetVerifiedCoins()

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollectionWithCallback(
			verified, coins.Resource{}, extendResourcesWithTradingVolumesCallback(explorer),
		),
	})
}

func extendResourcesWithTradingVolumesCallback(explorer *core.Explorer) func(resource.ParamInterface) resource.ParamsInterface {
	return func(model resource.ParamInterface) resource.ParamsInterface {
		return resource.ParamsInterface{
			extendResourceWithTradingVolumesParams(explorer, model.(models.Coin)),
		}
	}
}

func extendResourceWithTradingVolumesParams(explorer *core.Explorer, model models.Coin) coins.Params {
	return coins.Params{
		TradingVolume24h: explorer.CoinService.GetDailyTradingVolume(model.ID),
		TradingVolume1mo: explorer.CoinService.GetMonthlyTradingVolume(model.ID),
		PriceUsd:         explorer.PoolService.GetCoinPrice(uint64(model.ID)),
	}
}
