package transactions

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	transactions := r.Group("/transactions")
	{
		transactions.GET("", GetTransactions)
		transactions.GET("/:hash", GetTransaction)

	}
}
