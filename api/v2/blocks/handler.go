package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
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

// Blocks cache helpers
const CacheBlocksCount = 1

type CacheBlocksData struct {
	Blocks     []models.Block
	Pagination tools.Pagination
}

// Get list of blocks
func GetBlocks(c *gin.Context) {
	var blockModels []models.Block
	explorer := c.MustGet("explorer").(*core.Explorer)

	// fetch blocks
	pagination := tools.NewPagination(c.Request)

	getBlocks := func() []models.Block {
		return explorer.BlockRepository.GetPaginated(&pagination)
	}

	// cache last blocks
	if pagination.GetCurrentPage() == 1 && pagination.GetPerPage() == config.DefaultPaginationLimit {
		cached := explorer.Cache.Get("blocks", func() interface{} {
			return CacheBlocksData{getBlocks(), pagination}
		}, CacheBlocksCount).(CacheBlocksData)

		blockModels = cached.Blocks
		pagination = cached.Pagination
	} else {
		blockModels = getBlocks()
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
		"data": new(blocks.ResourceDetailed).Transform(*block),
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

	txs, err = explorer.TransactionService.PrepareTransactionsModel(txs)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination))
}
