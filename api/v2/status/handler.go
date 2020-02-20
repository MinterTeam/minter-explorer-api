package status

import (
	"errors"
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

const LastDataCacheTime = time.Duration(60)
const SlowAvgBlocksCacheTime = time.Duration(300)
const StatusPageCacheTime = 1

type Data struct {
	Result interface{}
	Error  error
}

func GetStatus(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// initialize channels for required status data
	txTotalCountCh := make(chan Data)
	avgTimeCh := make(chan Data)
	txTotalCount24hCh := make(chan Data)
	lastBlockCh := make(chan Data)

	go getTotalTxCountByLastDay(explorer, txTotalCount24hCh)
	go getTotalTxCount(explorer, txTotalCountCh)
	go getLastBlock(explorer, lastBlockCh)
	go getAverageBlockTime(explorer, avgTimeCh)

	// get values from goroutines
	txTotalCount, txTotalCount24h := <-txTotalCountCh, <-txTotalCount24hCh
	lastBlockData, avgBlockTime := <-lastBlockCh, <-avgTimeCh

	// handle errors from goroutines
	helpers.CheckErr(txTotalCount.Error)
	helpers.CheckErr(txTotalCount24h.Error)
	helpers.CheckErr(lastBlockData.Error)
	helpers.CheckErr(avgBlockTime.Error)

	// prepare data
	lastBlock := lastBlockData.Result.(models.Block)
	marketPrice := explorer.MarketService.PriceChange

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"bipPriceUsd":           marketPrice.Price,
			"bipPriceChange":        marketPrice.Change,
			"latestBlockHeight":     lastBlock.ID,
			"totalTransactions":     txTotalCount.Result.(int),
			"averageBlockTime":      avgBlockTime.Result.(float64),
			"latestBlockTime":       lastBlock.CreatedAt.Format(time.RFC3339),
			"marketCap":             getMarketCap(helpers.CalculateEmission(lastBlock.ID), marketPrice.Price),
			"transactionsPerSecond": getTransactionSpeed(txTotalCount24h.Result.(int)),
		},
	})
}

func GetStatusPage(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// initialize channels for required status data
	avgTimeCh := make(chan Data)
	txTotalCountCh := make(chan Data)
	slowBlocksTimeSumCh := make(chan Data)
	activeValidatorsCh := make(chan Data)
	activeCandidatesCh := make(chan Data)
	lastBlockCh := make(chan Data)
	tx24hDataCh := make(chan Data)
	stakesSumCh := make(chan Data)
	customCoinsDataCh := make(chan Data)

	go getTransactionsDataBy24h(explorer, tx24hDataCh)
	go getTotalTxCount(explorer, txTotalCountCh)
	go getActiveValidatorsCount(explorer, activeValidatorsCh)
	go getActiveCandidatesCount(explorer, activeCandidatesCh)
	go getAverageBlockTime(explorer, avgTimeCh)
	go getLastBlock(explorer, lastBlockCh)
	go getSumSlowBlocksTime(explorer, slowBlocksTimeSumCh)
	go getStakesSum(explorer, stakesSumCh)
	go getCustomCoinsData(explorer, customCoinsDataCh)

	// get values from goroutines
	tx24hData, txTotalCount := <-tx24hDataCh, <-txTotalCountCh
	activeValidators, activeCandidates := <-activeValidatorsCh, <-activeCandidatesCh
	avgBlockTime, lastBlockData, slowBlocksTimeSum := <-avgTimeCh, <-lastBlockCh, <-slowBlocksTimeSumCh
	stakesSumData, customCoinsData := <-stakesSumCh, <-customCoinsDataCh

	// handle errors from goroutines
	helpers.CheckErr(tx24hData.Error)
	helpers.CheckErr(txTotalCount.Error)
	helpers.CheckErr(activeValidators.Error)
	helpers.CheckErr(activeCandidates.Error)
	helpers.CheckErr(avgBlockTime.Error)
	helpers.CheckErr(lastBlockData.Error)
	helpers.CheckErr(slowBlocksTimeSum.Error)
	helpers.CheckErr(stakesSumData.Error)
	helpers.CheckErr(customCoinsData.Error)

	// prepare data
	lastBlock := lastBlockData.Result.(models.Block)
	tx24h := tx24hData.Result.(transaction.Tx24hData)
	customCoins := customCoinsData.Result.(coins.CustomCoinsStatusData)
	stakesSum := stakesSumData.Result.(string)

	status := "down"
	if isActive(lastBlock) {
		status = "active"
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"status":              status,
			"numberOfBlocks":      lastBlock.ID,
			"blockSpeed24h":       avgBlockTime.Result.(float64),
			"txTotalCount":        txTotalCount.Result.(int),
			"tx24hCount":          tx24hData.Result.(transaction.Tx24hData).Count,
			"activeValidators":    activeValidators.Result.(int),
			"activeCandidates":    activeCandidates.Result.(int),
			"totalDelegatedBip":   stakesSum,
			"customCoinsCount":    customCoins.Count,
			"averageTxCommission": helpers.Unit2Bip(tx24h.FeeAvg),
			"totalCommission":     helpers.Unit2Bip(tx24h.FeeSum),
			"customCoinsSum":      helpers.PipStr2Bip(customCoins.ReserveSum),
			"bipEmission":         helpers.CalculateEmission(lastBlock.ID),
			"freeFloatBip":        getFreeBipSum(stakesSum, lastBlock.ID),
			"txPerSecond":         getTransactionSpeed(tx24h.Count),
			"uptime":              calculateUptime(slowBlocksTimeSum.Result.(float64)),
		},
	})
}

func getTotalTxCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("total_tx_count"), func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(nil)
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getLastBlock(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("last_block", func() interface{} {
		return explorer.BlockRepository.GetLastBlock()
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getActiveCandidatesCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("active_candidates_count", func() interface{} {
		return explorer.ValidatorRepository.GetActiveCandidatesCount()
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getActiveValidatorsCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("active_validators_count", func() interface{} {
		return len(explorer.ValidatorRepository.GetActiveValidatorIds())
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getAverageBlockTime(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("avg_block_time", func() interface{} {
		return explorer.BlockRepository.GetAverageBlockTime()
	}, SlowAvgBlocksCacheTime)

	ch <- Data{data, nil}
}

func getSumSlowBlocksTime(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("slow_blocks_count", func() interface{} {
		return explorer.BlockRepository.GetSumSlowBlocksTimeBy24h()
	}, SlowAvgBlocksCacheTime)

	ch <- Data{data, nil}
}

func getTransactionsDataBy24h(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("tx_24h_data", func() interface{} {
		return explorer.TransactionRepository.Get24hTransactionsData()
	}, LastDataCacheTime).(transaction.Tx24hData)

	ch <- Data{data, nil}
}

func getTotalTxCountByLastDay(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	data := explorer.Cache.Get("last_day_total_tx_count", func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(&startTime)
	}, LastDataCacheTime)

	ch <- Data{data, nil}
}

func getStakesSum(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("stakes_sum"), func() interface{} {
		sum, err := explorer.StakeRepository.GetSumInBipValue()
		helpers.CheckErr(err)

		return helpers.PipStr2Bip(sum)
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getCustomCoinsData(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("custom_coins_data"), func() interface{} {
		data, err := explorer.CoinRepository.GetCustomCoinsStatusData()
		helpers.CheckErr(err)

		return data
	}, StatusPageCacheTime)

	ch <- Data{data, nil}
}

func getFreeBipSum(stakesSum string, lastBlockId uint64) float64 {
	stakes, err := strconv.ParseFloat(stakesSum, 64)
	helpers.CheckErr(err)

	return float64(helpers.CalculateEmission(lastBlockId)) - stakes
}

func getTransactionSpeed(total int) float64 {
	return helpers.Round(float64(total)/float64(86400), 8)
}

func calculateUptime(slow float64) float64 {
	return math.Round(((1-(slow/86400))*100)*100) / 100
}

func isActive(lastBlock models.Block) bool {
	return time.Now().Unix()-lastBlock.CreatedAt.Unix() <= config.NetworkActivePeriod
}

func getMarketCap(bipCount uint64, fiatPrice float64) float64 {
	return float64(bipCount) * fiatPrice
}

func recoveryStatusData(ch chan Data) {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			ch <- Data{nil, e}
		} else {
			ch <- Data{nil, errors.New(err.(string))}
		}
	}
}
