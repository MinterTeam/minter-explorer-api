package pools

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
)

type swapRouterResponse struct {
	Path   []string `json:"path"`
	Result string   `json:"result"`
}

type swapRouterErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type routeResource struct {
	SwapType  string               `json:"swap_type"`
	AmountIn  string               `json:"amount_in"`
	AmountOut string               `json:"amount_out"`
	Coins     []resource.Interface `json:"coins"`
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
		errors.SetErrorResponse(http.StatusNotFound, "Route path not exists.", c)
		return
	}

	var reqQuery FindSwapPoolRouteRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	coinFrom, coinTo, err := req.GetCoins(explorer)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Coins not found.", c)
		return
	}

	coin0 := strconv.FormatUint(uint64(coinFrom.ID), 10)
	coin1 := strconv.FormatUint(uint64(coinTo.ID), 10)

	resp, err := proxySwapPoolRouteRequest(coin0, coin1, reqQuery.Amount, reqQuery.TradeType)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Route path not exists.", c)
		return
	}

	if resp.IsError() {
		c.JSON(resp.StatusCode(), resp.Error())
		return
	}

	data := resp.Result().(*swapRouterResponse)
	path := make([]resource.Interface, len(data.Path))
	for i, cidStr := range data.Path {
		cid, _ := strconv.ParseUint(cidStr, 10, 64)
		coin, _ := explorer.CoinRepository.FindByID(uint(cid))
		path[i] = new(coins.IdResource).Transform(coin)
	}

	outputAmount := helpers.PipStr2Bip(data.Result)
	inputAmount := helpers.PipStr2Bip(reqQuery.Amount)

	if reqQuery.TradeType == "output" {
		inputAmount = helpers.PipStr2Bip(data.Result)
		outputAmount = helpers.PipStr2Bip(reqQuery.Amount)
	}

	c.JSON(resp.StatusCode(), routeResource{
		SwapType:  "pool",
		AmountIn:  inputAmount,
		AmountOut: outputAmount,
		Coins:     path,
	})
}

func proxySwapPoolRouteRequest(coin0, coin1, amount, tradeType string) (*resty.Response, error) {
	// todo: move host url to config
	return resty.New().R().
		SetError(&swapRouterErrorResponse{}).
		SetResult(&swapRouterResponse{}).
		Get(fmt.Sprintf("%s/v2/best_trade/%s/%s/%s/%s?max_depth=4", config.SwapRouterProxyUrl, coin0, coin1, tradeType, amount))
}
