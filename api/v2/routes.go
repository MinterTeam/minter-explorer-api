package apiV2

import (
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/addresses"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/statistics"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/status"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/transactions"
	"github.com/MinterTeam/minter-explorer-api/v2/api/v2/validators"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	v2 := r.Group("/v2")
	{
		blocks.ApplyRoutes(v2)
		coins.ApplyRoutes(v2)
		addresses.ApplyRoutes(v2)
		transactions.ApplyRoutes(v2)
		validators.ApplyRoutes(v2)
		statistics.ApplyRoutes(v2)
		status.ApplyRoutes(v2)
	}
}
