package transactions

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	blocks := r.Group("/transactions")
	{
		blocks.GET("", GetTransactions)
		blocks.GET("/:hash", GetTransaction)
	}
}
