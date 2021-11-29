package pools

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"os"
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

	resp, err := proxySwapPoolRouteRequest(req.Coin0, req.Coin1, reqQuery.Amount, reqQuery.TradeType)

	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	if resp.IsError() {
		c.JSON(resp.StatusCode(), resp.Error())
		return
	}

	c.JSON(resp.StatusCode(), resp.Result())
}

func proxySwapPoolRouteRequest(coin0, coin1, amount, tradeType string) (*resty.Response, error) {
	// todo: move host url to config
	hostUrl := "https://swap-router-api.minter.network"
	if os.Getenv("APP_BASE_COIN") == "MNT" {
		hostUrl = "https://swap-router-api.testnet.minter.network"
	}

	return resty.New().R().
		SetError(&swapRouterErrorResponse{}).
		SetResult(&swapRouterResponse{}).
		Get(fmt.Sprintf("%s/api/v1/pools/%s/%s/route?amount=%s&type=%s", hostUrl, coin0, coin1, amount, tradeType))
}
