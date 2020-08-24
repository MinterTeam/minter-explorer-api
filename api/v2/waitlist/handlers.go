package waitlist

import (
	"github.com/MinterTeam/minter-explorer-api/api/v2/addresses"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/waitlist"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetWaitlistByAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	minterAddress, err := addresses.GetAddressFromRequestUri(c)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	wl, err := explorer.WaitlistRepository.GetListByAddress(*minterAddress)
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, resource.TransformCollection(wl, waitlist.Resource{}))
}
