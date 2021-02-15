package pools

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	pools := r.Group("/pools")
	{
		pools.GET("", GetSwapPools)
		pools.GET("/token/:token", GetSwapPool)
		pools.GET("/token/:token/transactions", GetSwapPoolTransactions)
		pools.GET("/token/:token/providers", GetSwapPoolProviders)
		pools.GET("/token/:token/providers/:address", GetSwapPoolProvider)
		pools.GET("/coins/:coin0/:coin1", GetSwapPool)
		pools.GET("/coins/:coin0/:coin1/route", FindSwapPoolRoute)
		pools.GET("/coins/:coin0/:coin1/transactions", GetSwapPoolTransactions)
		pools.GET("/coins/:coin0/:coin1/providers", GetSwapPoolProviders)
		pools.GET("/coins/:coin0/:coin1/providers/:address", GetSwapPoolProvider)
		pools.GET("/autocomplete/coins", GetAutocompleteCoins)
		pools.GET("/providers/:address", GetSwapPoolsByProvider)
	}
}
