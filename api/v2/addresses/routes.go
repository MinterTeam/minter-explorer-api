package addresses

import "github.com/gin-gonic/gin"

// ApplyRoutes applies router to the gin Engine
func ApplyRoutes(r *gin.RouterGroup) {
	addresses := r.Group("/addresses")
	{
		addresses.GET("", GetAddresses)
		addresses.GET("/:address", GetAddress)
		addresses.GET("/:address/locks", GetLocks)
		addresses.GET("/:address/orders", GetLimitOrders)
		addresses.GET("/:address/transactions", GetTransactions)
		addresses.GET("/:address/events/unbonds", GetUnbonds)
		addresses.GET("/:address/events/slashes", GetSlashes)
		addresses.GET("/:address/events/bans", GetBans)
		addresses.GET("/:address/delegations", GetDelegations)
		addresses.GET("/:address/statistics/rewards", GetRewardsStatistics)
		addresses.GET("/:address/events/rewards/aggregated", GetAggregatedRewards)
	}
}
