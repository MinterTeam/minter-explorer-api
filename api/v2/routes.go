package apiV2

import (
	"github.com/MinterTeam/minter-explorer-api/api/v2/addresses"
	"github.com/MinterTeam/minter-explorer-api/api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/api/v2/statistics"
	"github.com/MinterTeam/minter-explorer-api/api/v2/status"
	"github.com/MinterTeam/minter-explorer-api/api/v2/transactions"
	"github.com/MinterTeam/minter-explorer-api/api/v2/validators"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v1 := r.Group("/v2")
	{
		blocks.ApplyRoutes(v1)
		coins.ApplyRoutes(v1)
		addresses.ApplyRoutes(v1)
		transactions.ApplyRoutes(v1)
		validators.ApplyRoutes(v1)
		statistics.ApplyRoutes(v1)
		status.ApplyRoutes(v1)
	}
}
