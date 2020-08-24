package waitlist

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	validators := r.Group("/waitlist")
	{
		validators.GET("/:address", GetWaitlistByAddress)
	}
}
