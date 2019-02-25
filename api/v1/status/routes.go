package status

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetStatus(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)
	startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")

	avgTimeCh := make(chan float64)
	lastBlockCh := make(chan models.Block)
	totalCountCh := make(chan int)
	totalCount24hCh := make(chan int)

	go func() { totalCount24hCh <- explorer.TransactionRepository.GetTotalTransactionCount(&startTime) }()
	go func() { lastBlockCh <- explorer.BlockRepository.GetLastBlock() }()
	go func() { totalCountCh <- explorer.TransactionRepository.GetTotalTransactionCount(nil) }()
	go func() { avgTimeCh <- explorer.BlockRepository.GetAverageBlockTime() }()

	txCount24h, lastBlock, txCountTotal, avgBlockTime := <-totalCount24hCh, <-lastBlockCh, <-totalCountCh, <-avgTimeCh

	c.JSON(http.StatusOK, gin.H{
		"bipPriceUsd":           0.07,                     //TODO: поменять значения, как станет ясно откуда брать
		"bipPriceBtc":           0.0000015883176063418346, //TODO: поменять значения, как станет ясно откуда брать
		"bipPriceChange":        10,                       //TODO: поменять значения, как станет ясно откуда брать
		"marketCap":             "",
		"latestBlockHeight":     lastBlock.ID,
		"latestBlockTime":       lastBlock.CreatedAt,
		"totalTransactions":     txCountTotal,
		"transactionsPerSecond": helpers.Round(float64(txCount24h)/float64(86400), 8),
		"averageBlockTime":      avgBlockTime,
	})
}

func GetStatusPage(c *gin.Context) {

}
