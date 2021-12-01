package addresses

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/address"
	"github.com/MinterTeam/minter-explorer-api/v2/aggregated_reward"
	"github.com/MinterTeam/minter-explorer-api/v2/chart"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/delegation"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/order"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/slash"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/unbond"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// FilterQueryRequest TODO: replace string to int
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

// GetAddresses Get list of addresses
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
		addresses[k] = extendModelWithBaseSymbolBalance(addr, addr.Address, explorer.Environment.Basecoin)
	}

	// extend the model array with empty model if not exists
	if len(addresses) != len(minterAddresses) {
		for _, item := range minterAddresses {
			if isModelsContainAddress(item, addresses) {
				continue
			}

			addresses = append(addresses, makeEmptyAddressModel(item, explorer.Environment.Basecoin))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(addresses, address.Resource{}),
	})
}

// GetAddress Get address detail
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

	addressModel := explorer.Cache.Get(fmt.Sprintf("address-%s-sum-%t", *minterAddress, request.WithSum), func() interface{} {
		balance, err := explorer.AddressService.GetBalance(*minterAddress, request.WithSum)
		helpers.CheckErr(err)

		if request.WithSum {
			return new(address.Resource).Transform(*balance.Model, address.Params{
				TotalBalanceSum:    balance.TotalBalanceSum,
				TotalBalanceSumUSD: balance.TotalBalanceSumUSD,
				StakeBalanceSum:    balance.StakeBalanceSum,
				StakeBalanceSumUSD: balance.StakeBalanceSumUSD,
			})
		}

		return new(address.Resource).Transform(*balance.Model)
	}, 1).(resource.Interface)

	c.JSON(http.StatusOK, gin.H{
		"data":              addressModel,
		"latest_block_time": explorer.Cache.GetLastBlock().Timestamp,
	})
}

// GetTransactions Get list of transactions by Minter address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var req AddressTransactionsQueryRequest
	err = c.ShouldBindQuery(&req)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	data := explorer.Cache.ExecuteOrGet(fmt.Sprintf("address-txs-%s", *minterAddress), func() interface{} {
		txs := explorer.TransactionRepository.GetPaginatedTxsByAddresses(
			[]string{*minterAddress},
			transaction.SelectFilter{
				SendType:   req.SendType,
				StartBlock: req.StartBlock,
				EndBlock:   req.EndBlock,
			}, &pagination)

		txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
		helpers.CheckErr(err)

		return resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination)
	}, 1, pagination.GetCurrentPage() != 1 || req.StartBlock != nil || req.EndBlock != nil || req.SendType != nil || pagination.Pager.GetLimit() != config.DefaultPaginationLimit).(resource.PaginationResource)

	c.JSON(http.StatusOK, data)
}

// GetAggregatedRewards Get aggregated by day list of address rewards
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

// GetSlashes
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

// GetDelegations
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

// GetRewardsStatistics Get rewards statistics by minter address
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

func GetBans(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	ids, err := explorer.StakeRepository.GetAddressValidatorIds(*minterAddress)
	if err != nil {
		log.WithError(err).Error("failed to load address validators ", *minterAddress)
		return
	}

	pagination := tools.NewPagination(c.Request)
	if len(ids) == 0 {
		c.JSON(http.StatusOK, resource.TransformPaginatedCollection([]models.ValidatorBan{}, events.AddressBanResource{}, pagination))
		return
	}

	bans, err := explorer.ValidatorRepository.GetBansByValidatorIds(ids, &pagination)
	if err != nil {
		log.WithError(err).Error("failed to load bans for address ", *minterAddress)
		return
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(bans, events.AddressBanResource{}, pagination))
}

type GetOrdersRequest struct {
	Status string `form:"status"  binding:"omitempty,oneof=canceled active expired partially_filled filled"`
}

// GetLimitOrders Get list of limit orders created by address
func GetLimitOrders(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request
	var rq GetOrdersRequest
	if err := c.ShouldBindQuery(&rq); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	orders, err := explorer.OrderRepository.GetListPaginated(&pagination, order.NewAddressFilter(*minterAddress), order.NewStatusFilter(rq.Status))
	if err != nil {
		log.WithError(err).Panic("failed to load orders")
		return
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(orders, new(order.Resource), pagination))
	return
}
