package checks

import (
	"github.com/MinterTeam/minter-explorer-api/v2/check"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetCheckRequest struct {
	RawCheck string `uri:"raw"`
}

type GetChecksRequest struct {
	FromAddress string `form:"address_to"     binding:"omitempty,minterAddress"`
	ToAddress   string `form:"address_from"   binding:"omitempty,minterAddress"`
}

// Get check information
func GetCheck(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetCheckRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	minterCheck, err := explorer.CheckRepository.GetByRawCheck(helpers.RemoveMinterPrefix(req.RawCheck))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Check not found.", c)
		return
	}

	checkData, err := explorer.TransactionService.TransformBaseCheckToModel(req.RawCheck)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, gin.H{
		"data": new(check.Resource).Transform(minterCheck, check.Params{CheckData: checkData}),
	})
}

// Get redeemed checks by filter
func GetChecks(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetChecksRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	filter := check.SelectFilter{FromAddress: req.FromAddress, ToAddress: req.ToAddress}
	checks, err := explorer.CheckRepository.GetListByFilter(filter, &pagination)
	helpers.CheckErr(err)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		checkData, err := explorer.TransactionService.TransformBaseCheckToModel(`Mc` + model.(models.Check).Data)
		helpers.CheckErr(err)
		return resource.ParamsInterface{check.Params{CheckData: checkData}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(checks, check.Resource{}, pagination, resourceCallback))
}
