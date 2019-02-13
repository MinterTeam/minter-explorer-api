package addresses

import (
	"github.com/MinterTeam/minter-explorer-api/address"
	"github.com/MinterTeam/minter-explorer-api/core"
	"github.com/MinterTeam/minter-explorer-api/errors"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GetAddressRequest struct {
	Address string `uri:"address" binding:"minterAddress"`
}

type GetAddressesRequest struct {
	Addresses []string `form:"addresses[]" binding:"required,minterAddress,max=50"`
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
	var minterAddresses []string
	for _, addr := range request.Addresses {
		minterAddresses = append(minterAddresses, helpers.RemoveMinterWalletPrefix(addr))
	}

	// fetch addresses
	addresses := explorer.AddressRepository.GetByAddresses(minterAddresses)

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(*addresses, address.Resource{}),
	})
}

// Get address detail
func GetAddress(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var request GetAddressRequest
	err := c.ShouldBindUri(&request)
	if err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// fetch address
	minterAddress := helpers.RemoveMinterWalletPrefix(request.Address)
	model := explorer.AddressRepository.GetByAddress(minterAddress)

	// if no models found
	if model == nil {
		model = &models.Address{
			Address: minterAddress,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(address.Resource).Transform(*model),
	})
}