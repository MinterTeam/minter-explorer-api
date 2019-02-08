package coins

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
)

// Get list of coins
func GetCoins(c *gin.Context) {
	db := c.MustGet(`db`).(*pg.DB)
	symbol := c.Query(`symbol`)

	// create service
	coinService := coins.CoinService{DB: db}

	// if symbol is not specified
	if symbol == "" {
		// fetch coins resource
		data := coinService.GetList()

		c.JSON(http.StatusOK, gin.H{"data": data})
	} else {
		// fetch coin
		coin := coinService.GetBySymbol(symbol)
		if coin == nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Coin not found.", c)
			return
		}

		data := []coins.CoinResource{*coin}

		c.JSON(http.StatusOK, gin.H{"data": data})
	}
}
