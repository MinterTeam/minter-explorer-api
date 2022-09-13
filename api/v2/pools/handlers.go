package pools

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/order"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-api/v2/transaction"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"math/big"
	"net/http"
	"strconv"
	"time"
)

type CachePoolsList struct {
	Pools resource.PaginationResource
}

const (
	CachePoolCoinsBlockCount = 1
	CachePoolsListBlockCount = 5
)

func GetSwapPool(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	data := explorer.Cache.Get(fmt.Sprintf("pools-%s-%s-%s", req.Coin0, req.Coin1, req.Token), func() interface{} {
		p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token})
		if err != nil {
			return nil
		}

		tv1d := explorer.PoolService.GetLastDayTradesVolume(p)
		tv30d := explorer.PoolService.GetLastMonthTradesVolume(p)

		return new(pool.Resource).Transform(p, pool.Params{
			FirstCoin:      req.Coin0,
			SecondCoin:     req.Coin1,
			TradeVolume1d:  tv1d.BipVolume,
			TradeVolume30d: tv30d.BipVolume,
		})
	}, 1)

	if data == nil {
		errors.SetErrorResponse(http.StatusNotFound, "Pool not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": data.(resource.Interface)})
}

func GetSwapPoolProvider(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolProviderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindProvider(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token}, helpers.RemoveMinterPrefix(req.Address))
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Provider not found.", c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": new(pool.ProviderResource).Transform(p, pool.Params{
			FirstCoin:  req.Coin0,
			SecondCoin: req.Coin1,
		}),
	})
}

// GetSwapPools get swap pool list
func GetSwapPools(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pools := explorer.Cache.ExecuteOrGet("pools", func() interface{} {
		return fetchPools(req, c)
	}, CachePoolsListBlockCount, len(c.Request.URL.Query()) != 0).(resource.PaginationResource)

	c.JSON(http.StatusOK, pools)
}

// fetch pools from db by request
func fetchPools(req GetSwapPoolsRequest, c *gin.Context) resource.PaginationResource {
	explorer := c.MustGet("explorer").(*core.Explorer)
	pagination := tools.NewPagination(c.Request)

	pools, _ := explorer.PoolRepository.GetPools(pool.SelectPoolsFilter{
		Coin:            req.Coin,
		ProviderAddress: req.Address,
	}, &pagination)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		p := model.(models.LiquidityPool)
		tv1d := explorer.PoolService.GetLastDayTradesVolume(p)
		tv30d := explorer.PoolService.GetLastMonthTradesVolume(p)

		return resource.ParamsInterface{pool.Params{
			TradeVolume1d:  tv1d.BipVolume,
			TradeVolume30d: tv30d.BipVolume,
		}}
	}

	return resource.TransformPaginatedCollectionWithCallback(pools, pool.Resource{}, pagination, resourceCallback)
}

func GetSwapPoolProviders(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolProvidersRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	providers, err := explorer.PoolRepository.GetProviders(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token}, &pagination)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Provider not found.", c)
		return
	}

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		return resource.ParamsInterface{pool.Params{
			FirstCoin:  req.Coin0,
			SecondCoin: req.Coin1,
		}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(providers, pool.ProviderResource{}, pagination, resourceCallback))
}

func GetSwapPoolsByProvider(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolsByProviderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)
	pools, err := explorer.PoolRepository.GetPoolsByProvider(helpers.RemoveMinterPrefix(req.Address), &pagination)
	helpers.CheckErr(err)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		return resource.ParamsInterface{pool.Params{}}
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollectionWithCallback(pools, pool.ProviderResource{}, pagination, resourceCallback))
}

func GetSwapPoolTransactions(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolTransactionsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	pagination := tools.NewPagination(c.Request)

	data := explorer.Cache.ExecuteOrGet(fmt.Sprintf("pools-txs-%s-%s-%s", req.Coin0, req.Coin1, req.Token), func() interface{} {
		txs := explorer.TransactionRepository.GetPaginatedTxsByFilter(transaction.PoolsFilter{
			Coin0: req.Coin0,
			Coin1: req.Coin1,
			Token: req.Token,
		}, &pagination)

		txs, err := explorer.TransactionService.PrepareTransactionsModel(txs)
		helpers.CheckErr(err)

		return resource.TransformPaginatedCollection(txs, transaction.Resource{}, pagination)
	}, 3, pagination.GetCurrentPage() != 1).(resource.PaginationResource)

	c.JSON(http.StatusOK, data)
}

func GetCoinsList(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	data := explorer.Cache.Get("pools_coins_list", func() interface{} {
		poolCoins, _ := explorer.PoolRepository.GetPoolsCoins()
		return resource.TransformCollection(poolCoins, coins.IdResource{})
	}, CachePoolCoinsBlockCount).([]resource.Interface)

	c.JSON(http.StatusOK, data)
}

func GetCoinPossibleSwaps(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetCoinPossibleSwapsRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	var reqQuery GetCoinPossibleSwapsRequestQuery
	if err := c.ShouldBindQuery(&reqQuery); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	poolCoins := explorer.Cache.Get("pools_coins_list_models", func() interface{} {
		poolCoins, _ := explorer.PoolRepository.GetPoolsCoins()
		return poolCoins
	}, CachePoolCoinsBlockCount).([]models.Coin)

	// find coin from by id or symbol
	var fromCoin *models.Coin
	if id, err := strconv.ParseUint(req.Coin, 10, 64); err == nil {
		for _, pc := range poolCoins {
			if uint64(pc.ID) == id {
				fromCoin = &pc
			}
		}
	} else {
		symbol, version := helpers.GetSymbolAndDefaultVersionFromStr(req.Coin)
		for _, pc := range poolCoins {
			if pc.Symbol == symbol && uint64(pc.Version) == version {
				fromCoin = &models.Coin{
					ID:     pc.ID,
					Symbol: pc.Symbol,
				}
			}
		}
	}

	// send empty response if coin not exists in pools
	if fromCoin == nil {
		c.JSON(http.StatusOK, resource.TransformCollection([]models.Coin{}, coins.IdResource{}))
		return
	}

	// TODO: hot-fix, remove
	c.JSON(http.StatusOK, resource.TransformCollection(poolCoins, coins.IdResource{}))
	return
}

func EstimateSwap(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req FindSwapPoolRouteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
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

	// define trade type
	tradeType := pool.TradeTypeExactInput
	if reqQuery.TradeType == "output" {
		tradeType = pool.TradeTypeExactOutput
	}

	coin0 := strconv.FormatUint(uint64(coinFrom.ID), 10)
	coin1 := strconv.FormatUint(uint64(coinTo.ID), 10)

	bancorAmount, bancorErr := explorer.SwapService.EstimateByBancor(coinFrom, coinTo, reqQuery.GetAmount(), tradeType)
	poolResp, poolErr := proxySwapPoolRouteRequest(coin0, coin1, reqQuery.Amount, reqQuery.TradeType)

	if poolErr != nil && bancorErr != nil && !poolResp.IsError() {
		errors.SetErrorResponse(http.StatusNotFound, "Route path not exists.", c)
		return
	}

	if bancorErr == nil && (poolErr != nil || poolResp.IsError()) {
		c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
		return
	}

	if bancorErr == nil {
		poolRespData := poolResp.Result().(*swapRouterResponse)
		outputAmount := helpers.StringToBigInt(poolRespData.Result)
		inputAmount := reqQuery.GetAmount()

		if tradeType == pool.TradeTypeExactOutput {
			inputAmount = helpers.StringToBigInt(poolRespData.Result)
			outputAmount = reqQuery.GetAmount()
		}

		if tradeType == pool.TradeTypeExactInput && bancorAmount.Cmp(outputAmount) >= 1 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}

		if tradeType == pool.TradeTypeExactOutput && bancorAmount.Cmp(inputAmount) <= 0 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}
	}

	if poolResp.IsError() {
		c.JSON(poolResp.StatusCode(), poolResp.Error())
		return
	}

	data := poolResp.Result().(*swapRouterResponse)
	path := make([]resource.Interface, len(data.Path))
	duplications := make(map[uint]bool)
	for i, cidStr := range data.Path {
		cid, _ := strconv.ParseUint(cidStr, 10, 64)
		coin, _ := explorer.CoinRepository.FindByID(uint(cid))
		path[i] = new(coins.IdResource).Transform(coin)

		// todo: temp logs; remove
		_, duplicationExists := duplications[coin.ID]
		if duplicationExists {
			log.Debugf("duplications: %v ; %d ; %v", data, cid, coin)
		} else {
			duplications[coin.ID] = true
		}
	}

	outputAmount := helpers.PipStr2Bip(data.Result)
	inputAmount := helpers.PipStr2Bip(reqQuery.Amount)

	if reqQuery.TradeType == "output" {
		inputAmount = helpers.PipStr2Bip(data.Result)
		outputAmount = helpers.PipStr2Bip(reqQuery.Amount)
	}

	c.JSON(poolResp.StatusCode(), routeResource{
		SwapType:  "pool",
		AmountIn:  inputAmount,
		AmountOut: outputAmount,
		Coins:     path,
	})
}

func GetSwapPoolTradesVolume(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request query
	var query GetSwapPoolTradesVolumeRequestQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token})
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Pool not found.", c)
		return
	}

	tradesVolume, err := explorer.PoolService.GetTradesVolume(p, query.Scale, nil)
	helpers.CheckErr(err)

	c.JSON(http.StatusOK, gin.H{
		"data": resource.TransformCollection(tradesVolume, pool.TradesVolumeResource{}),
	})
}

func GetSwapPoolsList(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	pools, err := explorer.PoolRepository.GetTracked()
	if err != nil {
		log.Panicf("failed to load pools: %s", err)
	}

	type cmcResource struct {
		BaseId              string `json:"base_id"`
		BaseName            string `json:"base_name"`
		BaseSymbol          string `json:"base_symbol"`
		BaseUsdPrice        string `json:"base_usd_price"`
		BaseChainId         string `json:"base_chain_id"`
		QuoteId             string `json:"quote_id"`
		QuoteName           string `json:"quote_name"`
		QuoteSymbol         string `json:"quote_symbol"`
		QuoteUsdPrice       string `json:"quote_usd_price"`
		QuoteChainId        string `json:"quote_chain_id"`
		LastPrice           string `json:"last_price"`
		BaseVolume          string `json:"base_volume"`
		QuoteVolume         string `json:"quote_volume"`
		BaseTradeVolume24h  string `json:"base_trade_volume_24h"`
		QuoteTradeVolume24h string `json:"quote_trade_volume_24h"`
		UsdTradeVolume24h   string `json:"usd_trade_volume_24h"`
	}

	resources := make(map[string]cmcResource, len(pools))

	coinToContractMap := make(map[uint64]*models.TokenContract)
	for _, p := range pools {
		if _, ok := coinToContractMap[p.FirstCoinId]; !ok {
			tc, _ := explorer.PoolRepository.GetTokenContractByCoinId(p.FirstCoinId)
			coinToContractMap[p.FirstCoinId] = tc
		}

		if _, ok := coinToContractMap[p.SecondCoinId]; !ok {
			tc, _ := explorer.PoolRepository.GetTokenContractByCoinId(p.SecondCoinId)
			coinToContractMap[p.SecondCoinId] = tc
		}
	}

	for _, p := range pools {
		startTime := time.Now().AddDate(0, 0, -1)
		trades, _ := explorer.PoolService.GetTradesVolume(p, nil, &startTime)

		tv := pool.TradeVolume{
			FirstCoinVolume:  "0",
			SecondCoinVolume: "0",
			BipVolume:        big.NewFloat(0),
		}

		if len(trades) != 0 {
			tv = trades[0]
		}

		usdTradeVolume24h := new(big.Float).Mul(
			explorer.PoolService.GetCoinPrice(0),
			tv.BipVolume,
		)

		baseContract, baseChain := helpers.GetTokenContractAndChain(coinToContractMap[p.FirstCoinId])
		quoteContract, quoteChain := helpers.GetTokenContractAndChain(coinToContractMap[p.SecondCoinId])

		ticker := fmt.Sprintf(`"%s_%s"`, baseContract, quoteContract)
		price := new(big.Float).Quo(
			helpers.StrToBigFloat(p.FirstCoinVolume),
			helpers.StrToBigFloat(p.SecondCoinVolume),
		)

		resources[ticker] = cmcResource{
			BaseId:              baseContract,
			BaseName:            p.FirstCoin.Name,
			BaseSymbol:          p.FirstCoin.GetSymbol(),
			BaseUsdPrice:        explorer.PoolService.GetCoinPrice(p.FirstCoinId).Text('f', 18),
			BaseChainId:         baseChain,
			BaseTradeVolume24h:  helpers.PipStr2Bip(tv.FirstCoinVolume),
			QuoteId:             quoteContract,
			QuoteName:           p.SecondCoin.Name,
			QuoteSymbol:         p.SecondCoin.GetSymbol(),
			QuoteUsdPrice:       explorer.PoolService.GetCoinPrice(p.SecondCoinId).Text('f', 18),
			QuoteTradeVolume24h: helpers.PipStr2Bip(tv.SecondCoinVolume),
			QuoteChainId:        quoteChain,
			LastPrice:           helpers.Bip2Str(price),
			BaseVolume:          helpers.PipStr2Bip(p.FirstCoinVolume),
			QuoteVolume:         helpers.PipStr2Bip(p.SecondCoinVolume),
			UsdTradeVolume24h:   usdTradeVolume24h.Text('f', 18),
		}
	}

	c.JSON(http.StatusOK, resources)
}

// GetAllSwapPools get all swap pool list
func GetAllSwapPools(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// add params to each model resource
	resourceCallback := func(model resource.ParamInterface) resource.ParamsInterface {
		p := model.(models.LiquidityPool)
		tv1d := explorer.PoolService.GetLastDayTradesVolume(p)
		tv30d := explorer.PoolService.GetLastMonthTradesVolume(p)

		return resource.ParamsInterface{pool.Params{
			TradeVolume1d:  tv1d.BipVolume,
			TradeVolume30d: tv30d.BipVolume,
		}}
	}

	c.JSON(http.StatusOK, resource.TransformCollectionWithCallback(explorer.PoolService.GetPools(), new(pool.Resource), resourceCallback))
}

// GetSwapPoolOrders Get orders related to liquidity pool
func GetSwapPoolOrders(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	// validate request
	var rq GetSwapPoolOrdersRequest
	if err := c.ShouldBindQuery(&rq); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	p, err := explorer.PoolRepository.FindByCoins(pool.SelectByCoinsFilter{Coin0: req.Coin0, Coin1: req.Coin1, Token: req.Token})
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Pool not found.", c)
		return
	}

	fromCoinId := uint64(0)
	if id, err := strconv.ParseUint(req.Coin0, 10, 64); err == nil {
		fromCoinId = id
	} else {
		fromCoinId, _ = explorer.CoinRepository.FindIdBySymbol(req.Coin0)
	}

	pagination := tools.NewPagination(c.Request)
	orders, err := explorer.OrderRepository.GetListPaginated(&pagination, order.NewPoolFilter(p),
		order.NewAddressFilter(helpers.RemoveMinterPrefix(rq.Address)), order.NewTypeFilter(rq.Type, p, fromCoinId),
		order.NewStatusFilter(rq.Status),
	)

	if err != nil {
		log.WithError(err).Panic("failed to load orders for pool")
		return
	}

	c.JSON(http.StatusOK, resource.TransformPaginatedCollection(orders, new(order.Resource), pagination))
}

// GetSwapPoolOrder Get limit order by id
func GetSwapPoolOrder(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// validate request
	var req GetSwapPoolOrderRequest
	if err := c.ShouldBindUri(&req); err != nil {
		errors.SetValidationErrorResponse(err, c)
		return
	}

	orderModel, err := explorer.OrderRepository.FindById(req.OrderId)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, "Order not found.", c)
		return
	}

	c.JSON(http.StatusOK, new(order.Resource).Transform(orderModel))
}
