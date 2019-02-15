package addresses

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	blocks := r.Group("/addresses")
	{
		blocks.GET("", GetAddresses)
		blocks.GET("/:address", GetAddress)
		blocks.GET("/:address/transactions", GetTransactions)
		blocks.GET("/:address/events/rewards", GetRewards)
		blocks.GET("/:address/events/slashes", GetSlashes)
	}
}
