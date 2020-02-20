package addresses

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/address"
	"github.com/MinterTeam/minter-explorer-api/aggregated_reward"
	"github.com/MinterTeam/minter-explorer-api/chart"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/delegation"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/events"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/slash"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/gin-gonic/gin"
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

	// extend the model array with empty model if not exists
	if len(addresses) != len(minterAddresses) {
		for _, item := range minterAddresses {
			if isModelsContainAddress(item, addresses) {
				continue
			}

			addresses = append(addresses, *makeEmptyAddressModel(item, explorer.Environment.BaseCoin))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(addresses, address.Resource{}),
	})
}

type NodeAddressBalanceResponse struct {
	Result struct {
		Balance map[string]string `json:"balance"`
	} `json:"result"`
}

// Get address detail
func GetAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	minterAddress, err := getAddressFromRequestUri(c)
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

	// if model not found
	if model == nil || len(model.Balances) == 0 {
		addressWithPrefix := `Mx` + *minterAddress
		httpClient := tools.NewHttpClient(explorer.Environment.NodeApiHost)
		resp := new(NodeAddressBalanceResponse)
		err := httpClient.Get(fmt.Sprintf("/address?address=%s", addressWithPrefix), resp)
		if err != nil {
			model = makeEmptyAddressModel(*minterAddress, explorer.Environment.BaseCoin)
		} else {
			var md []*models.Balance
			for coin, value := range resp.Result.Balance {
				md = append(md, &models.Balance{
					Value: value,
					Coin:  &models.Coin{Symbol: coin},
				})
			}

			model = &models.Address{
				Address:  *minterAddress,
				Balances: md,
			}
		}
	}

	// calculate overall address balance in base coin and fiat
	if request.WithSum {
		// compute available balance from address balances
		totalBalanceSum := explorer.BalanceService.GetTotalBalance(model)
		totalBalanceSumUSD := explorer.BalanceService.GetTotalBalanceInUSD(totalBalanceSum)

		c.JSON(http.StatusOK, gin.H{
			"data": new(address.Resource).Transform(*model, address.Params{
				TotalBalanceSum:    totalBalanceSum,
				TotalBalanceSumUSD: totalBalanceSumUSD,
			}),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{"data": new(address.Resource).Transform(*model)})
}

// Get list of transactions by Minter address
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := getAddressFromRequestUri(c)
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

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get aggregated by day list of address rewards
func GetAggregatedRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := getAddressFromRequestUri(c)
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

	minterAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)

	// get address stakes
	stakesCh := make(chan helpers.ChannelData)
	go func(ch chan helpers.ChannelData) {
		value, err := explorer.StakeRepository.GetPaginatedByAddress(*minterAddress, &pagination)
		ch <- helpers.NewChannelData(value, err)
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

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithAdditionalFields(
		delegationsData.Value,
		delegation.Resource{},
		pagination,
		additionalFields,
	))
}

// Get rewards statistics by minter address
func GetRewardsStatistics(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := getAddressFromRequestUri(c)
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

func prepareEventsRequest(c *gin.Context) (*events.SelectFilter, *tools.Pagination, error) {
	minterAddress, err := getAddressFromRequestUri(c)
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
func getAddressFromRequestUri(c *gin.Context) (*string, error) {
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
func isModelsContainAddress(minterAddress string, models []models.Address) bool {
	for _, item := range models {
		if item.Address == minterAddress {
			return true
		}
	}

	return false
}
