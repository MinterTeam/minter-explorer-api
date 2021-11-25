package pools

import (
	"fmt"
	"github.com/MinterTeam/explorer-sdk/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
)

type swapRouterResponse struct {
	SwapType  string `json:"swap_type"`
	AmountIn  string `json:"amount_in"`
	AmountOut string `json:"amount_out"`
	Coins     []struct {
		Id     int    `json:"id"`
		Symbol string `json:"symbol"`
	} `json:"coins"`
}

type swapRouterErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

func ProxySwapPoolRoute(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	if req.Coin0 == req.Coin1 {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	var reqQuery FindSwapPoolRouteRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	fromCoinId, toCoinId := uint64(0), uint64(0)
	if id, err := strconv.ParseUint(req.Coin0, 10, 64); err == nil {
		fromCoinId = id
	} else {
		fromCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin0)
		if err != nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
			return
		}
	}

	if id, err := strconv.ParseUint(req.Coin1, 10, 64); err == nil {
		toCoinId = id
	} else {
		toCoinId, err = explorer.CoinRepository.FindIdBySymbol(req.Coin1)
		if err != nil {
			errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
			return
		}
	}

	// define trade type
	tradeType := swap.TradeTypeExactInput
	if reqQuery.TradeType == "output" {
		tradeType = swap.TradeTypeExactOutput
	}

	trade, err := explorer.PoolService.FindSwapRoutePath(fromCoinId, toCoinId, tradeType, helpers.StringToBigInt(reqQuery.Amount))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		coin, err := explorer.CoinRepository.FindByID(uint(t.CoinID))
		helpers.CheckErr(err)
		path[i] = coin
	}

	c.JSON(http.StatusOK, new(pool.RouteResource).Transform(path, trade))
}

func proxySwapPoolRouteRequest(coin0, coin1, amount, tradeType string) (*resty.Response, error) {
	return resty.New().R().
		SetError(&swapRouterErrorResponse{}).
		SetResult(&swapRouterResponse{}).
		Get(fmt.Sprintf("https://swap-router-api.minter.network/api/v1/pools/%s/%s/route?amount=%s&type=%s", coin0, coin1, amount, tradeType))
}
