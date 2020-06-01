package status

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	r.GET("/status", GetStatus)
	r.GET("/status-page", GetStatusPage)
	r.GET("/info", GetInfo)
}
