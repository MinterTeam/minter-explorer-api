package addresses

import (
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
)

// GetAddressFromRequestUri Get minter address from current request uri
func GetAddressFromRequestUri(c *gin.Context) (*string, error) {
	var request GetAddressRequest
	if err := c.ShouldBindUri(&request); err != nil {
		return nil, err
	}

	minterAddress := helpers.RemoveMinterPrefix(request.Address)
	return &minterAddress, nil
}

// Return model address with zero base coin
func makeEmptyAddressModel(minterAddress string, baseCoin string) *models.Address {
	return &models.Address{
		Address: minterAddress,
		Balances: []*models.Balance{{
			Coin: &models.Coin{
				Symbol: baseCoin,
				Type:   models.CoinTypeBase,
			},
			Value: "0",
		}},
	}
}

// Check that array of address models contain exact minter address
func isModelsContainAddress(minterAddress string, models []*models.Address) bool {
	for _, item := range models {
		if item.Address == minterAddress {
			return true
		}
	}

	return false
}

func extendModelWithBaseSymbolBalance(model *models.Address, minterAddress, baseCoin string) *models.Address {
	// if model not found
	if model == nil || len(model.Balances) == 0 {
		return makeEmptyAddressModel(minterAddress, baseCoin)
	}

	isBaseSymbolExists := false
	for _, b := range model.Balances {
		if b.CoinID == 0 {
			isBaseSymbolExists = true
		}
	}

	if !isBaseSymbolExists {
		model.Balances = append(model.Balances, &models.Balance{
			Value: "0",
			Coin:  &models.Coin{Symbol: baseCoin, Type: models.CoinTypeBase},
		})
	}

	return model
}

func prepareEventsRequest(c *gin.Context) (*events.SelectFilter, *tools.Pagination, error) {
	minterAddress, err := GetAddressFromRequestUri(c)
	if err != nil {
		return nil, nil, err
	}

	var requestQuery FilterQueryRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		return nil, nil, err
	}

	pagination := tools.NewPagination(c.Request)

	return &events.SelectFilter{
		Address:    *minterAddress,
		StartBlock: requestQuery.StartBlock,
		EndBlock:   requestQuery.EndBlock,
	}, &pagination, nil
}
