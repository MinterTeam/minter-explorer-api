package apiV1

import (
	"github.com/MinterTeam/minter-explorer-api/api/v1/blocks"
	"github.com/MinterTeam/minter-explorer-api/api/v1/coins"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		blocks.ApplyRoutes(v1)
		coins.ApplyRoutes(v1)
	}
}
