package statistics

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	statistics := r.Group("/statistics")
	{
		statistics.GET("/transactions", GetTransactions)
	}
}
