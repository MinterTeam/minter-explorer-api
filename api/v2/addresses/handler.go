package addresses

import (
	"github.com/MinterTeam/minter-explorer-api/v2/address"
	"github.com/MinterTeam/minter-explorer-api/v2/aggregated_reward"
	"github.com/MinterTeam/minter-explorer-api/v2/chart"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/delegation"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/slash"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/unbond"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type GetAddressRequest struct {
	Address string `uri:"address" binding:"minterAddress"`
}

type GetAddressRequestQuery struct {
	WithSum bool `form:"with_sum"`
}

type GetAddressesRequest struct {
	Addresses []string `form:"addresses[]" binding:"required,minterAddress,max=50"`
}

// TODO: replace string to int
type FilterQueryRequest struct {
	StartBlock *string `form:"start_block" binding:"omitempty,numeric"`
	EndBlock   *string `form:"end_block"   binding:"omitempty,numeric"`
	Page       *string `form:"page"        binding:"omitempty,numeric"`
}

type AddressTransactionsQueryRequest struct {
	SendType   *string `form:"send_type"   binding:"omitempty,oneof=incoming outcoming"`
	StartBlock *string `form:"start_block" binding:"omitempty,numeric"`
	EndBlock   *string `form:"end_block"   binding:"omitempty,numeric"`
	Page       *string `form:"page"        binding:"omitempty,numeric"`
}

type StatisticsQueryRequest struct {
	StartTime *string `form:"start_time" binding:"omitempty,timestamp"`
	EndTime   *string `form:"end_time"   binding:"omitempty,timestamp"`
}

type AggregatedRewardsQueryRequest struct {
	StartTime *string `form:"start_time" binding:"omitempty,timestamp"`
	EndTime   *string `form:"end_time"   binding:"omitempty,timestamp"`
}

// Get list of addresses
func GetAddresses(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetAddressesRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// remove Minter wallet prefix from each address
	minterAddresses := make([]string, len(request.Addresses))
	for key, addr := range request.Addresses {
		minterAddresses[key] = helpers.RemoveMinterPrefix(addr)
	}

	// fetch addresses
	addresses := explorer.AddressRepository.GetByAddresses(minterAddresses)

	for k, addr := range addresses {
		addresses[k] = extendModelWithBaseSymbolBalance(addr, addr.Address, explorer.Environment.BaseCoin)
	}

	// extend the model array with empty model if not exists
	if len(addresses) != len(minterAddresses) {
		for _, item := range minterAddresses {
			if isModelsContainAddress(item, addresses) {
				continue
			}

			addresses = append(addresses, makeEmptyAddressModel(item, explorer.Environment.BaseCoin))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(addresses, address.Resource{}),
	})
}

// Get address detail
func GetAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query params
	var request GetAddressRequestQuery
	if err := c.ShouldBindQuery(&request); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch address
	model := explorer.AddressRepository.GetByAddress(*minterAddress)
	model = extendModelWithBaseSymbolBalance(model, *minterAddress, explorer.Environment.BaseCoin)

	// calculate overall address balance in base coin and fiat
	if request.WithSum {
		stakes, err := explorer.StakeRepository.GetAllByAddress(model.Address)
		if err != nil {
			log.Fatal(err)
		}

		// compute available balance from address balances
		totalBalanceSum := explorer.BalanceService.GetTotalBalance(model)
		totalBalanceSumUSD := explorer.BalanceService.GetTotalBalanceInUSD(totalBalanceSum)
		stakeBalanceSum := explorer.BalanceService.GetStakeBalance(stakes)
		stakeBalanceSumUSD := explorer.BalanceService.GetTotalBalanceInUSD(stakeBalanceSum)

		c.JSON(http.StatusOK, gin.H{
			"data": new(address.Resource).Transform(*model, address.Params{
				TotalBalanceSum:    totalBalanceSum,
				TotalBalanceSumUSD: totalBalanceSumUSD,
				StakeBalanceSum:    stakeBalanceSum,
				StakeBalanceSumUSD: stakeBalanceSumUSD,
			}),
			"latest_block_time": explorer.Cache.GetLastBlock().Timestamp,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":              new(address.Resource).Transform(*model),
		"latest_block_time": explorer.Cache.GetLastBlock().Timestamp,
	})
}

// Get list of transactions by Minter address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var requestQuery AddressTransactionsQueryRequest
	err = c.ShouldBindQuery(&requestQuery)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByAddresses(
		[]string{*minterAddress},
		transaction.SelectFilter{
			SendType:   requestQuery.SendType,
			StartBlock: requestQuery.StartBlock,
			EndBlock:   requestQuery.EndBlock,
		}, &pagination)

	txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get aggregated by day list of address rewards
func GetAggregatedRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery AggregatedRewardsQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	rewards := explorer.RewardRepository.GetPaginatedAggregatedByAddress(aggregated_reward.SelectFilter{
		Address:   *minterAddress,
		StartTime: requestQuery.StartTime,
		EndTime:   requestQuery.EndTime,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(rewards, aggregated_reward.Resource{}, pagination))
}

// Get list of slashes by Minter address
func GetSlashes(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	slashes := explorer.SlashRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(slashes, slash.Resource{}, *pagination))
}

// Get list of delegations by Minter address
func GetDelegations(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// get address stakes
	stakesCh := make(chan helpers.ChannelData)
	go func(ch chan helpers.ChannelData) {
		value, err := explorer.StakeRepository.GetAllByAddress(*minterAddress)
		if err != nil {
			ch <- helpers.NewChannelData(value, err)
			return
		}

		stakes, err := explorer.StakeService.PrepareStakesModels(value)
		ch <- helpers.NewChannelData(stakes, err)
	}(stakesCh)

	// get address total delegated sum in base coin
	stakesSumCh := make(chan helpers.ChannelData)
	go func(ch chan helpers.ChannelData) {
		value, err := explorer.StakeRepository.GetSumInBipValueByAddress(*minterAddress)
		ch <- helpers.NewChannelData(value, err)
	}(stakesSumCh)

	delegationsData, stakesSumData := <-stakesCh, <-stakesSumCh
	helpers.CheckErr(delegationsData.Error)
	helpers.CheckErr(stakesSumData.Error)

	additionalFields := map[string]interface{}{
		"total_delegated_bip_value": helpers.PipStr2Bip(
			stakesSumData.Value.(string),
		),
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(delegationsData.Value, delegation.Resource{}),
		"meta": gin.H{
			"additional": additionalFields,
		},
	})
}

// Get rewards statistics by minter address
func GetRewardsStatistics(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery StatisticsQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	chartData := explorer.RewardRepository.GetAggregatedChartData(aggregated_reward.SelectFilter{
		Address:   *minterAddress,
		EndTime:   requestQuery.EndTime,
		StartTime: requestQuery.StartTime,
	})

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(chartData, chart.RewardResource{}),
	})
}

func GetUnbonds(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	wl, err := explorer.UnbondRepository.GetListByAddress(filter, pagination)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(wl, unbond.Resource{}),
	})
}

func prepareEventsRequest(c *gin.Context) (*events.SelectFilter, *tools.Pagination, error) {
	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		return nil, nil, err
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		return nil, nil, err
	}

	pagination := tools.NewPagination(c.Request)

	return &events.SelectFilter{
		Address:    *minterAddress,
		StartBlock: requestQuery.StartBlock,
		EndBlock:   requestQuery.EndBlock,
	}, &pagination, nil
}

// Get minter address from current request uri
func GetAddressFromRequestUri(c *gin.Context) (*string, error) {
	var request GetAddressRequest
	if err := c.ShouldBindUri(&request); err != nil {
		return nil, err
	}

	minterAddress := helpers.RemoveMinterPrefix(request.Address)
	return &minterAddress, nil
}

// Return model address with zero base coin
func makeEmptyAddressModel(minterAddress string, baseCoin string) *models.Address {
	return &models.Address{
		Address: minterAddress,
		Balances: []*models.Balance{{
			Coin: &models.Coin{
				Symbol: baseCoin,
			},
			Value: "0",
		}},
	}
}

// Check that array of address models contain exact minter address
func isModelsContainAddress(minterAddress string, models []*models.Address) bool {
	for _, item := range models {
		if item.Address == minterAddress {
			return true
		}
	}

	return false
}

func extendModelWithBaseSymbolBalance(model *models.Address, minterAddress, baseCoin string) *models.Address {
	// if model not found
	if model == nil || len(model.Balances) == 0 {
		return makeEmptyAddressModel(minterAddress, baseCoin)
	}

	isBaseSymbolExists := false
	for _, b := range model.Balances {
		if b.CoinID == 0 {
			isBaseSymbolExists = true
		}
	}

	if !isBaseSymbolExists {
		model.Balances = append(model.Balances, &models.Balance{
			Value: "0",
			Coin:  &models.Coin{Symbol: baseCoin},
		})
	}

	return model
}
