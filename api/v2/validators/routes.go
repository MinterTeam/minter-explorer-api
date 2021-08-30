package validators

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	validators := r.Group("/validators")
	{
		validators.GET("", GetValidators)
		validators.GET("/meta", GetValidatorsMeta)
		validators.GET("/:publicKey", GetValidator)
		validators.GET("/:publicKey/stakes", GetValidatorStakes)
		validators.GET("/:publicKey/transactions", GetValidatorTransactions)
		validators.GET("/:publicKey/events/bans", GetValidatorBans)
		validators.GET("/:publicKey/events/slashes", GetValidatorSlashes)
	}
}
