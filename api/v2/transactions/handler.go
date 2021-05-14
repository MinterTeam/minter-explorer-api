package transactions

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/invalid_transaction"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO: replace string in StartBlock, EndBlock, Page to int
type GetTransactionsRequest struct {
	Addresses  []string `form:"addresses[]"  binding:"omitempty,minterAddress"`
	Page       string   `form:"page"         binding:"omitempty,numeric"`
	StartBlock *string  `form:"start_block"  binding:"omitempty,numeric"`
	EndBlock   *string  `form:"end_block"    binding:"omitempty,numeric"`
}

type GetTransactionRequest struct {
	Hash string `uri:"hash" binding:"minterTxHash"`
}

// Transaction cache helpers
const CacheBlocksCount = 1

type CacheTxData struct {
	Transactions []models.Transaction
	Pagination   tools.Pagination
}

// Get list of transactions
func GetTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionsRequest
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

	// fetch data
	pagination := tools.NewPagination(c.Request)

	var txs []models.Transaction
	if len(minterAddresses) > 0 {
		txs = explorer.TransactionRepository.GetPaginatedTxsByAddresses(minterAddresses, transaction.SelectFilter{
			StartBlock: request.StartBlock,
			EndBlock:   request.EndBlock,
		}, &pagination)

		txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
		helpers.CheckErr(err)
	} else {
		// prepare retrieving models
		getTxsFunc := func() []models.Transaction {
			txs :=  explorer.TransactionRepository.GetPaginatedTxsByFilter(blocks.RangeSelectFilter{
				StartBlock: request.StartBlock,
				EndBlock:   request.EndBlock,
			}, &pagination)

			txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
			helpers.CheckErr(err)

			return txs
		}

		// cache last transactions
		if len(c.Request.URL.Query()) == 0 {
			cached := explorer.Cache.Get("transactions", func() interface{} {
				return CacheTxData{getTxsFunc(), pagination}
			}, CacheBlocksCount).(CacheTxData)

			txs = cached.Transactions
			pagination = cached.Pagination
		} else {
			txs = getTxsFunc()
		}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}

// Get transaction detail by hash
func GetTransaction(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetTransactionRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch data
	hash := helpers.RemoveMinterPrefix(request.Hash)
	tx := explorer.TransactionRepository.GetTxByHash(hash)
	if tx == nil {
		invalidTx := explorer.InvalidTransactionRepository.GetTxByHash(hash)
		if invalidTx == nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Transaction not found.", c)
			return
		}

		c.JSON(http.StatusPartialContent, gin.H{
			"data": new(invalid_transaction.Resource).Transform(*invalidTx),
		})
		return
	}

	tx, err = explorer.TransactionService.PrepareTransactionModel(tx)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, gin.H{
		"data": new(transaction.Resource).Transform(*tx),
	})
}
