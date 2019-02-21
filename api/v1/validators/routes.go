package validators

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	validators := r.Group("/validators")
	{
		validators.GET("/:publicKey/transactions", GetValidatorTransactions)
		validators.GET("/:publicKey", GetValidator)
	}
}
