package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/paginator"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/gin-gonic/gin"
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
	explorer := c.MustGet("explorer").(*core.Explorer)

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
	models := explorer.BlockRepository.GetPaginated(page, CountOfBlocksPerPage)

	// make response as empty array if no models
	if len(models) == 0 {
		empty := make([]blocks.Resource, 0)
		c.JSON(http.StatusOK, gin.H{"data": empty})
		return
	}

	// transform to resource
	blocksList := resource.TransformCollection(models, blocks.Resource{})

	response := paginator.Resource{
		Data: blocksList,
		Links: paginator.LinksResource{
			First: "",
			Last:  "",
			Prev:  "",
			Next:  "",
		},
		Meta: paginator.MetaResource{
			CurrentPage: 1,
			From:        2,
			LastPage:    3,
			Path:        "",
			PerPage:     5,
			To:          6,
			Total:       7,
		},
	}

	c.JSON(http.StatusOK, response)
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

	// transform to resource
	var blocksResource blocks.Resource
	data := blocksResource.Transform(*block)

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
