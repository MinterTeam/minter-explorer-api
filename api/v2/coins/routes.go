package coins

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	coins := r.Group("/coins")
	{
		coins.GET("", GetCoins)
		coins.GET("symbol/:symbol", GetCoinBySymbol)
		coins.GET("id/:id", GetCoinByID)
		coins.GET("oracle/verified", GetOracleVerifiedCoins)
	}
}
