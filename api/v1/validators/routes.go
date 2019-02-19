package validators

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	blocks := r.Group("/validators")
	{
		blocks.GET("/:publicKey/transactions", GetValidatorTransactions)
	}
}
