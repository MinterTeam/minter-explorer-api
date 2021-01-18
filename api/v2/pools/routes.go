package pools

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	pools := r.Group("/pools")
	{
		pools.GET("/", GetSwapPools)
		pools.GET("/:coin0/:coin1", GetSwapPool)
		pools.GET("/:coin0/:coin1/:address", GetSwapPoolProvider)
	}
}
