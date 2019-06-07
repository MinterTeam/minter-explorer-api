package status

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const LastDataDataCacheTime = time.Duration(60)
const SlowAvgBlocksCacheTime = time.Duration(300)
const StatusPageCacheTime = 1

func GetStatus(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	totalCountCh := make(chan int)
	avgTimeCh := make(chan float64)
	totalCount24hCh := make(chan int)
	lastBlockCh := make(chan models.Block)

	go getTotalTxCountByLastDay(explorer, totalCount24hCh)
	go getTotalTxCount(explorer, totalCountCh)
	go getLastBlock(explorer, lastBlockCh)
	go getAverageBlockTime(explorer, avgTimeCh)

	txCount24h, lastBlock, txCountTotal, avgBlockTime := <-totalCount24hCh, <-lastBlockCh, <-totalCountCh, <-avgTimeCh
	// TODO: replace with tools.GetCurrentFiatPrice(explorer.Environment.BaseCoin, "USD")
	price := explorer.Environment.StartPriceUSD

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"bipPriceUsd":           price,                    //TODO: поменять значения, как станет ясно откуда брать
			"bipPriceChange":        10,                       //TODO: поменять значения, как станет ясно откуда брать
			"marketCap":             getMarketCap(helpers.CalculateEmission(lastBlock.ID), price),
			"latestBlockHeight":     lastBlock.ID,
			"latestBlockTime":       lastBlock.CreatedAt.Format(time.RFC3339),
			"totalTransactions":     txCountTotal,
			"transactionsPerSecond": getTransactionSpeed(txCount24h),
			"averageBlockTime":      avgBlockTime,
		},
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
	stakesSumCh := make(chan string)
	customCoinsDataCh := make(chan coins.CustomCoinsStatusData)

	go getTransactionsDataBy24h(explorer, tx24hDataCh)
	go getTotalTxCount(explorer, txTotalCountCh)
	go getActiveValidatorsCount(explorer, activeValidatorsCh)
	go getActiveCandidatesCount(explorer, activeCandidatesCh)
	go getAverageBlockTime(explorer, avgTimeCh)
	go getLastBlock(explorer, lastBlockCh)
	go getSlowBlocksCount(explorer, slowBlocksCountCh)
	go getStakesSum(explorer, stakesSumCh)
	go getCustomCoinsData(explorer, customCoinsDataCh)

	tx24hData, txTotalCount := <-tx24hDataCh, <-txTotalCountCh
	activeValidators, activeCandidates := <-activeValidatorsCh, <-activeCandidatesCh
	avgBlockTime, lastBlock, slowBlocksCount := <-avgTimeCh, <-lastBlockCh, <-slowBlocksCountCh
	stakesSum, customCoinsData := <-stakesSumCh, <-customCoinsDataCh

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
			"averageTxCommission": helpers.Unit2Bip(tx24hData.FeeAvg),
			"totalCommission":     helpers.Unit2Bip(tx24hData.FeeSum),
			"totalDelegatedBip":   stakesSum,
			"customCoinsSum":      helpers.PipStr2Bip(customCoinsData.ReserveSum),
			"customCoinsCount":    customCoinsData.Count,
			"freeFloatBip":        getFreeBipSum(stakesSum, lastBlock.ID),
			"bipEmission":         helpers.CalculateEmission(lastBlock.ID),
		},
	})
}

func getTotalTxCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.Cache.Get(fmt.Sprintf("total_tx_count"), func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(nil)
	}, StatusPageCacheTime).(int)
}

func getLastBlock(explorer *core.Explorer, ch chan models.Block) {
	ch <- explorer.Cache.Get("last_block", func() interface{} {
		return explorer.BlockRepository.GetLastBlock()
	}, StatusPageCacheTime).(models.Block)
}

func getActiveCandidatesCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.Cache.Get("active_candidates_count", func() interface{} {
		return explorer.ValidatorRepository.GetActiveCandidatesCount()
	}, StatusPageCacheTime).(int)
}

func getActiveValidatorsCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.Cache.Get("active_validators_count", func() interface{} {
		return len(explorer.ValidatorRepository.GetActiveValidatorIds())
	}, StatusPageCacheTime).(int)
}

func getAverageBlockTime(explorer *core.Explorer, ch chan float64) {
	ch <- explorer.Cache.Get("avg_block_time", func() interface{} {
		return explorer.BlockRepository.GetAverageBlockTime()
	}, SlowAvgBlocksCacheTime).(float64)
}

func getSlowBlocksCount(explorer *core.Explorer, ch chan int) {
	ch <- explorer.Cache.Get("slow_blocks_count", func() interface{} {
		return explorer.BlockRepository.GetSlowBlocksCountBy24h()
	}, SlowAvgBlocksCacheTime).(int)
}

func getTransactionsDataBy24h(explorer *core.Explorer, ch chan transaction.Tx24hData) {
	ch <- explorer.Cache.Get("tx_24h_data", func() interface{} {
		return explorer.TransactionRepository.Get24hTransactionsData()
	}, LastDataDataCacheTime).(transaction.Tx24hData)
}

func getTotalTxCountByLastDay(explorer *core.Explorer, ch chan int) {
	startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	ch <- explorer.Cache.Get("last_day_total_tx_count", func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(&startTime)
	}, LastDataDataCacheTime).(int)
}

func getStakesSum(explorer *core.Explorer, ch chan string) {
	ch <- explorer.Cache.Get(fmt.Sprintf("stakes_sum"), func() interface{} {
		sum, err := explorer.StakeRepository.GetSumInBipValue()
		helpers.CheckErr(err)

		return helpers.PipStr2Bip(sum)
	}, StatusPageCacheTime).(string)
}

func getCustomCoinsData(explorer *core.Explorer, ch chan coins.CustomCoinsStatusData) {
	ch <- explorer.Cache.Get(fmt.Sprintf("custom_coins_data"), func() interface{} {
		data, err := explorer.CoinRepository.GetCustomCoinsStatusData()
		helpers.CheckErr(err)

		return data
	}, StatusPageCacheTime).(coins.CustomCoinsStatusData)
}

func getFreeBipSum(stakesSum string, lastBlockId uint64) float64 {
	stakes, err := strconv.ParseFloat(stakesSum, 64)
	helpers.CheckErr(err)
	return float64(helpers.CalculateEmission(lastBlockId)) - stakes
}

func getTransactionSpeed(total int) float64 {
	return helpers.Round(float64(total)/float64(86400), 8)
}

func calculateUptime(total int, slow int) float64 {
	if total == 0 {
		return 0
	}

	return float64((1 - float64(slow)/float64(total)) * 100)
}

func isActive(lastBlock models.Block) bool {
	return time.Now().Unix()-lastBlock.CreatedAt.Unix() <= config.NetworkActivePeriod
}

func getMarketCap(bipCount uint64, fiatPrice float64) float64 {
	return float64(bipCount) * fiatPrice
}