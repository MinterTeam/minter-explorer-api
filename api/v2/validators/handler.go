package validators

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/slash"
	"github.com/MinterTeam/minter-explorer-api/v2/stake"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetValidatorRequest struct {
	PublicKey string `uri:"publicKey"    binding:"required,minterPubKey"`
}

// TODO: replace string to int
type GetValidatorTransactionsRequest struct {
	Page       string  `form:"page"         binding:"omitempty,numeric"`
	StartBlock *string `form:"start_block"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"end_block"    binding:"omitempty,numeric"`
}

// CacheBlocksCount cache time
const CacheBlocksCount = 3

// Get list of transaction by validator public key
func GetValidatorTransactions(c *gin.Context) {
	var validatorRequest GetValidatorRequest
	var request GetValidatorTransactionsRequest

	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	err := c.ShouldBindUri(&validatorRequest)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	err = c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	publicKey := helpers.RemoveMinterPrefix(validatorRequest.PublicKey)
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.ValidatorFilter{
		ValidatorPubKey: publicKey,
		StartBlock:      request.StartBlock,
		EndBlock:        request.EndBlock,
	}, &pagination)

	txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get validator detail by public key
func GetValidator(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetValidatorRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	data := explorer.ValidatorRepository.GetByPublicKey(helpers.RemoveMinterPrefix(request.PublicKey))

	// check validator to existing
	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	activeValidatorIDs := getActiveValidatorIDs(explorer)
	totalStake := getTotalStakeByActiveValidators(explorer, activeValidatorIDs)
	minStake := explorer.ValidatorService.GetMinStakesByValidator(data)

	c.JSON(http.StatusOK, gin.H{
		"data": validator.ResourceDetailed{}.Transform(*data, validator.Params{
			TotalStake:          totalStake,
			MinStake:            minStake,
			ActiveValidatorsIDs: activeValidatorIDs,
		}),
	})
}

// Get list of validators
func GetValidators(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// fetch validators
	validators := explorer.Cache.Get("validators", func() interface{} {
		validators := explorer.ValidatorRepository.GetValidatorsAndStakes()

		activeValidatorIDs := getActiveValidatorIDs(explorer)
		totalStake := getTotalStakeByActiveValidators(explorer, activeValidatorIDs)

		// add params to each model resource
		resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
			v := model.(models.Validator)
			minStake := explorer.ValidatorService.GetMinStakesByValidator(&v)

			params := validator.Params{
				MinStake:            minStake,
				TotalStake:          totalStake,
				ActiveValidatorsIDs: activeValidatorIDs,
			}

			return resource.ParamsInterface{params}
		}

		return resource.TransformCollectionWithCallback(
			validators,
			validator.ResourceDetailed{},
			resourceCallback,
		)
	}, CacheBlocksCount).([]resource.Interface)

	c.JSON(http.StatusOK, gin.H{"data": validators})
}

type GetValidatorDelegationsRequestQuery struct {
	Page string `form:"page" binding:"omitempty,numeric"`
}

// Get validator stake list
func GetValidatorStakes(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	request := new(GetValidatorRequest)
	err := c.ShouldBindUri(request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	requestQuery := new(GetValidatorDelegationsRequestQuery)
	if err := c.ShouldBindQuery(requestQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch validator by public key
	publicKey := helpers.RemoveMinterPrefix(request.PublicKey)
	v := explorer.ValidatorRepository.GetByPublicKey(publicKey)
	if v == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	// fetch validator stakes
	pagination := tools.NewPagination(c.Request)
	stakes, err := explorer.StakeRepository.GetPaginatedByValidator(*v, &pagination)
	helpers.CheckErr(err)

	stakes, err = explorer.StakeService.PrepareStakesModels(stakes)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(stakes, stake.Resource{}, pagination))
}

// Get validator slashes list
func GetValidatorSlashes(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	request := new(GetValidatorRequest)
	err := c.ShouldBindUri(request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch validator by public key
	publicKey := helpers.RemoveMinterPrefix(request.PublicKey)
	v := explorer.ValidatorRepository.GetByPublicKey(publicKey)
	if v == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	// fetch validator stakes
	pagination := tools.NewPagination(c.Request)
	slashes, err := explorer.SlashRepository.GetPaginatedByValidator(v, &pagination)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(slashes, slash.Resource{}, pagination))
}

// Get validator bans list
func GetValidatorBans(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	request := new(GetValidatorRequest)
	err := c.ShouldBindUri(request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch validator by public key
	publicKey := helpers.RemoveMinterPrefix(request.PublicKey)
	v := explorer.ValidatorRepository.GetByPublicKey(publicKey)
	if v == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Validator not found.", c)
		return
	}

	// fetch validator stakes
	pagination := tools.NewPagination(c.Request)
	bans, err := explorer.ValidatorRepository.GetBans(v, &pagination)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(bans, events.BanResource{}, pagination))
}

// GetValidatorsMeta Get list of validators meta
func GetValidatorsMeta(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// fetch validators
	validators := explorer.Cache.Get("validators-meta", func() interface{} {
		validators := explorer.ValidatorRepository.GetValidators()
		return resource.TransformCollection(validators, validator.Resource{})
	}, 120).([]resource.Interface)

	c.JSON(http.StatusOK, gin.H{"data": validators})
}
