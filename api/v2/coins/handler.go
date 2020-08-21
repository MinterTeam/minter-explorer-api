package coins

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const CacheBlocksCount = 1

// Get list of coins
func GetCoins(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)
	symbol := c.Query("symbol")

	var data []models.Coin

	if symbol == "" {
		// fetch coins resource
		data = explorer.Cache.Get("coins", func() interface{} {
			return explorer.CoinRepository.GetCoins()
		}, CacheBlocksCount).([]models.Coin)
	} else {
		// fetch coins by symbol
		data = explorer.CoinRepository.GetBySymbol(symbol)
	}

	// make response as empty array if no models found
	if len(data) == 0 {
		empty := make([]coins.Resource, 0)

		c.JSON(http.StatusOK, gin.H{"data": empty})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resource.TransformCollection(data, coins.Resource{})})
}
