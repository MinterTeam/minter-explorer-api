package addresses

import (
	"github.com/MinterTeam/minter-explorer-api/address"
	"github.com/MinterTeam/minter-explorer-api/aggregated_reward"
	"github.com/MinterTeam/minter-explorer-api/chart"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/MinterTeam/minter-explorer-api/delegation"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/events"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/reward"
	"github.com/MinterTeam/minter-explorer-api/slash"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
)

type GetAddressRequest struct {
	Address string `uri:"address" binding:"minterAddress"`
}

type GetAddressRequestQuery struct {
	WithSum bool `form:"withSum"`
}

type GetAddressesRequest struct {
	Addresses []string `form:"addresses[]" binding:"required,minterAddress,max=50"`
}

// TODO: replace string to int
type FilterQueryRequest struct {
	StartBlock *string `form:"startblock" binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"   binding:"omitempty,numeric"`
	Page       *string `form:"page"       binding:"omitempty,numeric"`
}

type StatisticsQueryRequest struct {
	Scale     *string `form:"scale"     binding:"omitempty,eq=minute|eq=hour|eq=day"`
	StartTime *string `form:"startTime" binding:"omitempty,timestamp"`
	EndTime   *string `form:"endTime"   binding:"omitempty,timestamp"`
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
	if model == nil {
		model = makeEmptyAddressModel(*minterAddress, explorer.Environment.BaseCoin)
	}

	// calculate overall address balance (with stakes) in base coin and fiat
	var availableBalanceSum, totalBalanceSum,
		totalBalanceSumUSD, availableBalanceSumUSD *big.Float
	if request.WithSum {
		// get address stakes
		addressStakes := explorer.StakeRepository.GetByAddress(*minterAddress)

		// compute available balance from address balances
		availableBalanceSum = explorer.BalanceService.GetAvailableBalance(model.Balances)
		availableBalanceSumUSD = explorer.BalanceService.GetBalanceSumInUSD(availableBalanceSum)

		// compute total balance from address balances and stakes
		totalBalanceSum = explorer.BalanceService.GetTotalBalance(model.Balances, addressStakes)
		totalBalanceSumUSD = explorer.BalanceService.GetBalanceSumInUSD(totalBalanceSum)

		c.JSON(http.StatusOK, gin.H{
			"data": new(address.Resource).Transform(*model, address.Params{
				AvailableBalanceSum:    availableBalanceSum,
				AvailableBalanceSumUSD: availableBalanceSumUSD,
				TotalBalanceSum:        totalBalanceSum,
				TotalBalanceSumUSD:     totalBalanceSumUSD,
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
	var requestQuery FilterQueryRequest
	err = c.ShouldBindQuery(&requestQuery)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByAddresses(
		[]string{*minterAddress},
		transaction.BlocksRangeSelectFilter{
			StartBlock: requestQuery.StartBlock,
			EndBlock:   requestQuery.EndBlock,
		}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get list of rewards by Minter address
func GetRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	filter, pagination, err := prepareEventsRequest(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	rewards := explorer.RewardRepository.GetPaginatedByAddress(*filter, pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(rewards, reward.Resource{}, *pagination))
}

func GetAggregatedRewards(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := getAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	pagination := tools.NewPagination(c.Request)
	rewards := explorer.RewardRepository.GetPaginatedAggregatedByAddress(reward.AggregatedSelectFilter{
		Address:    *minterAddress,
		StartBlock: requestQuery.StartBlock,
		EndBlock:   requestQuery.EndBlock,
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

	// fetch data
	pagination := tools.NewPagination(c.Request)
	delegations := explorer.StakeRepository.GetPaginatedByAddress(*minterAddress, &pagination)

	// fetch total delegated sum
	totalDelegated, err := explorer.StakeRepository.GetSumInBipValueByAddress(*minterAddress)
	helpers.CheckErr(err)

	// create additional field
	additionalFields := map[string]interface{}{"total_delegated_bip_value": helpers.PipStr2Bip(totalDelegated)}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithAdditionalFields(
		delegations, delegation.Resource{}, pagination, additionalFields))
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

	// set scale instead of default if exists
	scale := config.DefaultStatisticsScale
	if requestQuery.Scale != nil {
		scale = *requestQuery.Scale
	}

	// fetch data
	chartData := explorer.RewardRepository.GetChartData(*minterAddress, chart.SelectFilter{
		Scale:     scale,
		StartTime: requestQuery.StartTime,
		EndTime:   requestQuery.EndTime,
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
