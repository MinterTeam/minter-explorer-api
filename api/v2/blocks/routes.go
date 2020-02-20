package blocks

import (
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	blocks := r.Group("/blocks")
	{
		blocks.GET("", GetBlocks)
		blocks.GET("/:height", GetBlock)
		blocks.GET("/:height/transactions", GetBlockTransactions)
	}
}
