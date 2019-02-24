package apiV1

import (
	"github.com/MinterTeam/minter-explorer-api/api/v1/addresses"
	"github.com/MinterTeam/minter-explorer-api/api/v1/blocks"
	"github.com/MinterTeam/minter-explorer-api/api/v1/coins"
	"github.com/MinterTeam/minter-explorer-api/api/v1/statistics"
	"github.com/MinterTeam/minter-explorer-api/api/v1/transactions"
	"github.com/MinterTeam/minter-explorer-api/api/v1/validators"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v1")
	{
		blocks.ApplyRoutes(v1)
		coins.ApplyRoutes(v1)
		addresses.ApplyRoutes(v1)
		transactions.ApplyRoutes(v1)
		validators.ApplyRoutes(v1)
		statistics.ApplyRoutes(v1)
	}
}
