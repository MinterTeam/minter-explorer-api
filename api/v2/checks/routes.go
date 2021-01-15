package checks

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	blocks := r.Group("/checks")
	{
		blocks.GET("/", GetChecks)
		blocks.GET("/:raw", GetCheck)
	}
}
