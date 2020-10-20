package status

import (
	"errors"
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	lastDataCacheTime      = time.Duration(60)
	slowAvgBlocksCacheTime = time.Duration(300)
	pageCacheTime          = 1
	chainID                = "minter-mainnet-2"
)

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
			"bip_price_usd":           marketPrice.Price,
			"bip_price_change":        marketPrice.Change,
			"latest_block_height":     lastBlock.ID,
			"total_transactions":      txTotalCount.Result.(int),
			"avg_block_time":          avgBlockTime.Result.(float64),
			"latest_block_time":       lastBlock.CreatedAt.Format(time.RFC3339),
			"market_cap":              getMarketCap(helpers.CalculateEmission(lastBlock.ID), marketPrice.Price),
			"transactions_per_second": getTransactionSpeed(txTotalCount24h.Result.(int)),
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
			"status":                     status,
			"blocks_count":               lastBlock.ID,
			"block_speed_24h":            avgBlockTime.Result.(float64),
			"transactions_total":         txTotalCount.Result.(int),
			"transaction_count_24h":      tx24hData.Result.(transaction.Tx24hData).Count,
			"active_validators":          activeValidators.Result.(int),
			"active_candidates":          activeCandidates.Result.(int),
			"total_delegated_bip":        stakesSum,
			"custom_coins_count":         customCoins.Count,
			"avg_transaction_commission": helpers.Unit2Bip(tx24h.FeeAvg),
			"total_commission":           helpers.Unit2Bip(tx24h.FeeSum),
			"custom_coins_sum":           helpers.PipStr2Bip(customCoins.ReserveSum),
			"bip_emission":               helpers.CalculateEmission(lastBlock.ID),
			"free_float_bip":             getFreeBipSum(stakesSum, lastBlock.ID),
			"transactions_per_second":    getTransactionSpeed(tx24h.Count),
			"uptime":                     calculateUptime(slowBlocksTimeSum.Result.(float64)),
		},
	})
}

type InfoResponse struct {
	ChainId                string    `json:"chain_id"`
	BlockHeight            uint64    `json:"block_height"`
	BlockTime              float64   `json:"block_time"`
	TotalCirculatingTokens float64   `json:"total_circulating_tokens"`
	NotBondedTokens        float64   `json:"not_bonded_tokens"`
	BondedTokens           float64   `json:"bonded_tokens"`
	TotalTxsNum            int       `json:"total_txs_num"`
	TotalValidatorNum      int       `json:"total_validator_num"`
	Time                   time.Time `json:"time"`
}

func GetInfo(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// initialize channels for required status data
	avgTimeCh := make(chan Data)
	txTotalCountCh := make(chan Data)
	activeValidatorsCh := make(chan Data)
	lastBlockCh := make(chan Data)
	stakesSumCh := make(chan Data)
	customCoinsDataCh := make(chan Data)

	go getTotalTxCount(explorer, txTotalCountCh)
	go getActiveValidatorsCount(explorer, activeValidatorsCh)
	go getAverageBlockTime(explorer, avgTimeCh)
	go getLastBlock(explorer, lastBlockCh)
	go getStakesSum(explorer, stakesSumCh)
	go getCustomCoinsData(explorer, customCoinsDataCh)

	// get values from goroutines
	txTotalCount := <-txTotalCountCh
	activeValidators := <-activeValidatorsCh
	avgBlockTime, lastBlockData := <-avgTimeCh, <-lastBlockCh
	stakesSumData := <-stakesSumCh

	// handle errors from goroutines
	helpers.CheckErr(txTotalCount.Error)
	helpers.CheckErr(activeValidators.Error)
	helpers.CheckErr(avgBlockTime.Error)
	helpers.CheckErr(lastBlockData.Error)
	helpers.CheckErr(stakesSumData.Error)

	// prepare data
	lastBlock := lastBlockData.Result.(models.Block)
	stakesSum := stakesSumData.Result.(string)
	notBondedTokens := getFreeBipSum(stakesSum, lastBlock.ID)
	bondedTokens, err := strconv.ParseFloat(stakesSum, 64)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, InfoResponse{
		ChainId:                chainID,
		BlockHeight:            lastBlock.ID,
		BlockTime:              avgBlockTime.Result.(float64),
		TotalCirculatingTokens: bondedTokens + notBondedTokens,
		NotBondedTokens:        notBondedTokens,
		BondedTokens:           bondedTokens,
		TotalTxsNum:            txTotalCount.Result.(int),
		TotalValidatorNum:      activeValidators.Result.(int),
		Time:                   lastBlock.CreatedAt,
	})
}

func getTotalTxCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("total_tx_count"), func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(nil)
	}, pageCacheTime)

	ch <- Data{data, nil}
}

func getLastBlock(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("last_block", func() interface{} {
		return explorer.BlockRepository.GetLastBlock()
	}, pageCacheTime)

	ch <- Data{data, nil}
}

func getActiveCandidatesCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("active_candidates_count", func() interface{} {
		return explorer.ValidatorRepository.GetActiveCandidatesCount()
	}, pageCacheTime)

	ch <- Data{data, nil}
}

func getActiveValidatorsCount(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("active_validators_count", func() interface{} {
		return len(explorer.ValidatorRepository.GetActiveValidatorIds())
	}, pageCacheTime)

	ch <- Data{data, nil}
}

func getAverageBlockTime(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("avg_block_time", func() interface{} {
		return explorer.BlockRepository.GetAverageBlockTime()
	}, slowAvgBlocksCacheTime)

	ch <- Data{data, nil}
}

func getSumSlowBlocksTime(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("slow_blocks_count", func() interface{} {
		return explorer.BlockRepository.GetSumSlowBlocksTimeBy24h()
	}, slowAvgBlocksCacheTime)

	ch <- Data{data, nil}
}

func getTransactionsDataBy24h(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get("tx_24h_data", func() interface{} {
		return explorer.TransactionRepository.Get24hTransactionsData()
	}, lastDataCacheTime).(transaction.Tx24hData)

	ch <- Data{data, nil}
}

func getTotalTxCountByLastDay(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	startTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")
	data := explorer.Cache.Get("last_day_total_tx_count", func() interface{} {
		return explorer.TransactionRepository.GetTotalTransactionCount(&startTime)
	}, lastDataCacheTime)

	ch <- Data{data, nil}
}

func getStakesSum(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("stakes_sum"), func() interface{} {
		sum, err := explorer.StakeRepository.GetSumInBipValue()
		helpers.CheckErr(err)

		return helpers.PipStr2Bip(sum)
	}, pageCacheTime)

	ch <- Data{data, nil}
}

func getCustomCoinsData(explorer *core.Explorer, ch chan Data) {
	defer recoveryStatusData(ch)

	data := explorer.Cache.Get(fmt.Sprintf("custom_coins_data"), func() interface{} {
		data, err := explorer.CoinRepository.GetCustomCoinsStatusData()
		helpers.CheckErr(err)

		return data
	}, pageCacheTime)

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
