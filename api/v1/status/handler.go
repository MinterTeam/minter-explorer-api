package status

import (
	"github.com/MinterTeam/minter-explorer-api/core"
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

	go func() { tx24hDataCh <- explorer.TransactionRepository.Get24hTransactionsData() }()
	go func() { txTotalCountCh <- explorer.TransactionRepository.GetTotalTransactionCount(nil) }()
	go func() { activeValidatorsCh <- len(explorer.ValidatorRepository.GetActiveValidatorIds()) }()
	go func() { activeCandidatesCh <- explorer.ValidatorRepository.GetActiveCandidatesCount() }()
	go func() { avgTimeCh <- explorer.BlockRepository.GetAverageBlockTime() }()
	go func() { lastBlockCh <- explorer.BlockRepository.GetLastBlock() }()
	go func() { slowBlocksCountCh <- explorer.BlockRepository.GetSlowBlocksCountBy24h() }()

	tx24hData, txTotalCount := <-tx24hDataCh, <-txTotalCountCh
	activeValidators, activeCandidates := <-activeValidatorsCh, <-activeCandidatesCh
	avgBlockTime, lastBlock, slowBlocksCount := <-avgTimeCh, <-lastBlockCh, <-slowBlocksCountCh

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":              "active",
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

func getTransactionSpeed(total int) float64 {
	return helpers.Round(float64(total)/float64(86400), 8)
}

func calculateUptime(total int, slow int) float64 {
	if total == 0 {
		return 0
	}

	return float64((1 - float64(slow) / float64(total)) * 100)
}

func getStatus(lastBlock models.Block) bool {
	
}