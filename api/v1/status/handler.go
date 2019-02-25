package status

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetStatus(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)
	startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")

	totalCountCh := make(chan int)
	avgTimeCh := make(chan float64)
	totalCount24hCh := make(chan int)
	lastBlockCh := make(chan models.Block)

	go getTotalTxCount(explorer, totalCount24hCh, &startTime)
	go getTotalTxCount(explorer, totalCountCh, nil)
	go getLastBlock(explorer, lastBlockCh)
	go getAverageBlockTime(explorer, avgTimeCh)

	txCount24h, lastBlock, txCountTotal, avgBlockTime := <-totalCount24hCh, <-lastBlockCh, <-totalCountCh, <-avgTimeCh

	c.JSON(http.StatusOK, gin.H{
		"bipPriceUsd":           0.07,                     //TODO: поменять значения, как станет ясно откуда брать
		"bipPriceBtc":           0.0000015883176063418346, //TODO: поменять значения, как станет ясно откуда брать
		"bipPriceChange":        10,                       //TODO: поменять значения, как станет ясно откуда брать
		"marketCap":             "",
		"latestBlockHeight":     lastBlock.ID,
		"latestBlockTime":       lastBlock.CreatedAt,
		"totalTransactions":     txCountTotal,
		"transactionsPerSecond": getTransactionSpeed(txCount24h),
		"averageBlockTime":      avgBlockTime,
	})
}

func GetStatusPage(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	avgTimeCh := make(chan float64)
	txTotalCountCh := make(chan int)
	slowBlocksCountCh := make(chan int)
	activeValidatorsCh := make(chan int)
	activeCandidatesCh := make(chan int)
	lastBlockCh := make(chan models.Block)
	tx24hDataCh := make(chan transaction.Tx24hData)

	go getTransactionsDataBy24h(explorer, tx24hDataCh)
	go getTotalTxCount(explorer, txTotalCountCh, nil)
	go getActiveValidatorsCount(explorer, activeValidatorsCh)
	go getActiveCandidatesCount(explorer, activeCandidatesCh)
	go getAverageBlockTime(explorer, avgTimeCh)
	go getLastBlock(explorer, lastBlockCh)
	go getSlowBlocksCount(explorer, slowBlocksCountCh)

	tx24hData, txTotalCount := <-tx24hDataCh, <-txTotalCountCh
	activeValidators, activeCandidates := <-activeValidatorsCh, <-activeCandidatesCh
	avgBlockTime, lastBlock, slowBlocksCount := <-avgTimeCh, <-lastBlockCh, <-slowBlocksCountCh

	status := "down"
	if isActive(lastBlock) {
		status = "active"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":              status,
			"uptime":              calculateUptime(tx24hData.Count, slowBlocksCount),
			"numberOfBlocks":      lastBlock.ID,
			"blockSpeed24h":       avgBlockTime,
			"txTotalCount":        txTotalCount,
			"tx24hCount":          tx24hData.Count,
			"txPerSecond":         getTransactionSpeed(tx24hData.Count),
			"activeValidators":    activeValidators,
			"activeCandidates":    activeCandidates,
			"averageTxCommission": tx24hData.FeeAvg,
			"totalCommission":     tx24hData.FeeSum,
		},
	})
}

func getTotalTxCount(explorer *core.Explorer, ch chan int, startTime *string) {
	ch <- explorer.TransactionRepository.GetTotalTransactionCount(startTime)
}

func getLastBlock(explorer *core.Explorer, ch chan models.Block) {
	ch <- explorer.BlockRepository.GetLastBlock()
}

func getAverageBlockTime(explorer *core.Explorer, ch chan float64) {
	ch <- explorer.BlockRepository.GetAverageBlockTime()
}

func getActiveCandidatesCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.ValidatorRepository.GetActiveCandidatesCount()
}

func getActiveValidatorsCount(explorer *core.Explorer, ch chan int) {
	ch <- len(explorer.ValidatorRepository.GetActiveValidatorIds())
}

func getSlowBlocksCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.BlockRepository.GetSlowBlocksCountBy24h()
}

func getTransactionsDataBy24h(explorer *core.Explorer, ch chan transaction.Tx24hData) {
	ch <- explorer.TransactionRepository.Get24hTransactionsData()
}

func getTransactionSpeed(total int) float64 {
	return helpers.Round(float64(total)/float64(86400), 8)
}

func calculateUptime(total int, slow int) float64 {
	if total == 0 {
		return 0
	}

	return float64((1 - float64(slow) / float64(total)) * 100)
}

func isActive(lastBlock models.Block) bool {
	return time.Now().Unix() - lastBlock.CreatedAt.Unix() <= config.NetworkActivePeriod
}