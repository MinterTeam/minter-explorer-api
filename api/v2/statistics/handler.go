package statistics

import (
	"github.com/MinterTeam/minter-explorer-api/v2/chart"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type GetTransactionsRequest struct {
	Scale     *string `form:"scale"      binding:"omitempty,eq=minute|eq=hour|eq=day"`
	StartTime *string `form:"start_time" binding:"omitempty,timestamp"`
	EndTime   *string `form:"end_time"   binding:"omitempty,timestamp"`
}

// statistics cache time
const CacheTime = time.Duration(600)

func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	var request GetTransactionsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// set default scale
	scale := config.DefaultStatisticsScale
	if request.Scale != nil {
		scale = *request.Scale
	}

	// set default start time
	startTime := helpers.StartOfTheDay(time.Now().AddDate(0, 0, config.DefaultStatisticsDayDelta)).Format("2006-01-02 15:04:05")
	if request.StartTime != nil {
		startTime = *request.StartTime
	}

	txFunc := func() interface{} {
		return explorer.TransactionRepository.GetTxCountChartDataByFilter(chart.SelectFilter{
			Scale:     scale,
			StartTime: &startTime,
			EndTime:   request.EndTime,
		})
	}

	// cache request without query parameters
	var txs interface{}
	if len(c.Request.URL.Query()) == 0 {
		txs = explorer.Cache.Get("tx_statistics", txFunc, CacheTime)
	} else {
		txs = txFunc()
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(txs, chart.TransactionResource{}),
	})
}
