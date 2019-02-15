package validators

import (
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetValidatorRequest struct {
	PublicKey string `uri:"publicKey"    binding:"required,minterPubKey"`
}

type GetValidatorTransactionsRequest struct {
	Page       string  `form:"page"        binding:"omitempty,numeric"`
	StartBlock *string `form:"startblock"  binding:"omitempty,numeric"`
	EndBlock   *string `form:"endblock"    binding:"omitempty,numeric"`
}

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
	txs := explorer.TransactionRepository.GetPaginatedTxByFilter(transaction.SelectFilter{
		ValidatorPubKey: &publicKey,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}