package pagination

import (
	"github.com/gin-gonic/gin"
)

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("paginator", NewService(c.Request))
		c.Next()
	}
}
