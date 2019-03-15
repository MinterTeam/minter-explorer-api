package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-api/transaction"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TODO: replace string to int
type GetBlockRequest struct {
	ID string `uri:"height" binding:"numeric"`
}

// TODO: replace string to int
type GetBlocksRequest struct {
	Page string `form:"page" binding:"omitempty,numeric"`
}

const CacheBlocksCount = 1

// Get list of blocks
func GetBlocks(c *gin.Context) {
	var blockModels []models.Block
	explorer := c.MustGet("explorer").(*core.Explorer)

	// fetch blocks
	pagination := tools.NewPagination(c.Request)

	// cache last blocks
	if pagination.GetCurrentPage() == 1 {
		blockModels = explorer.Cache.Get("blocks", func() interface{} {
			return explorer.BlockRepository.GetPaginated(&pagination)
		}, CacheBlocksCount).([]models.Block)
	} else {
		blockModels = explorer.BlockRepository.GetPaginated(&pagination)
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(blockModels, blocks.Resource{}, pagination))
}

// Get block detail
func GetBlock(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetBlockRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// parse to uint64
	blockId, err := strconv.ParseUint(request.ID, 10, 64)
	helpers.CheckErr(err)

	// fetch block by height
	block := explorer.BlockRepository.GetById(blockId)

	// check block to existing
	if block == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Block not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(blocks.Resource).Transform(*block),
	})
}

// Get list of transactions by block height
func GetBlockTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetBlockRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var requestQuery GetBlocksRequest
	err = c.ShouldBindQuery(&requestQuery)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// parse to uint64
	blockId, err := strconv.ParseUint(request.ID, 10, 64)
	helpers.CheckErr(err)

	// fetch data
	pagination := tools.NewPagination(c.Request)
	txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.BlockFilter{
		BlockId: blockId,
	}, &pagination)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}
