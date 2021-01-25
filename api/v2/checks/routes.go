package checks

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	checks := r.Group("/checks")
	{
		checks.GET("", GetChecks)
		checks.GET("/:raw", GetCheck)
	}
}
