package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"net/http"
	"strconv"
)

type GetBlockRequest struct {
	ID string `uri:"height" binding:"numeric"`
}

type GetBlocksRequest struct {
	Page string `form:"page" binding:"omitempty,numeric"`
}

// Count of blocks per page
const CountOfBlocksPerPage = 50

// Get list of blocks
func GetBlocks(c *gin.Context) {
	db := c.MustGet(`db`).(*pg.DB)

	// validate request
	var request GetBlocksRequest
	err := c.ShouldBindQuery(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// set current page
	var page = 1
	if request.Page != "" {
		page, err = strconv.Atoi(request.Page)
		helpers.CheckErr(err)
	}

	// fetch blocks
	blockService := blocks.BlockService{DB: db}
	blocksList := blockService.GetList(page, CountOfBlocksPerPage)

	c.JSON(http.StatusOK, gin.H{
		"data": *blocksList,
	})
}

// Get block detail
func GetBlock(c *gin.Context) {
	db := c.MustGet(`db`).(*pg.DB)

	// validate request
	var request GetBlockRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// parse to uint64
	blockId, err := strconv.ParseUint(request.ID, 10, 64)

	// fetch block by height
	blockService := blocks.BlockService{DB: db}
	block := blockService.GetById(blockId)

	// check block to existing
	if block == nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Block not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": block,
	})
}
