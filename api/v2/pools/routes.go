package pools

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	pools := r.Group("/pools")
	{
		pools.GET("", GetSwapPools)
		pools.GET("/all", GetAllSwapPools)
		pools.GET("/token/:token", GetSwapPool)
		pools.GET("/token/:token/transactions", GetSwapPoolTransactions)
		pools.GET("/token/:token/providers", GetSwapPoolProviders)
		pools.GET("/token/:token/providers/:address", GetSwapPoolProvider)
		pools.GET("/token/:token/stats/volume", GetSwapPoolTradesVolume)
		pools.GET("/coins/:coin0/:coin1", GetSwapPool)
		pools.GET("/coins/:coin0/:coin1/route", ProxySwapPoolRoute)
		pools.GET("/coins/:coin0/:coin1/estimate", EstimateSwap)
		pools.GET("/coins/:coin0/:coin1/transactions", GetSwapPoolTransactions)
		pools.GET("/coins/:coin0/:coin1/providers", GetSwapPoolProviders)
		pools.GET("/coins/:coin0/:coin1/stats/volume", GetSwapPoolTradesVolume)
		pools.GET("/coins/:coin0/:coin1/providers/:address", GetSwapPoolProvider)
		pools.GET("/providers/:address", GetSwapPoolsByProvider)
		pools.GET("/list/coins", GetCoinsList)
		pools.GET("/list/coins/:coin", GetCoinPossibleSwaps)
		pools.GET("/list/cmc", GetSwapPoolsList)
	}
}
