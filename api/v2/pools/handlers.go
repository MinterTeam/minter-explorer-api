package pools

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/MinterTeam/explorer-sdk/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core"
	"github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
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
	"sync"
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
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pool not found.", c)
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
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
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
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Provider not found.", c)
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

func FindSwapPoolRoute(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// todo: remove, temp stats collector
	hasher := md5.New()
	hasher.Write([]byte(c.Request.RequestURI + time.Now().String() + c.ClientIP()))
	rid := hex.EncodeToString(hasher.Sum(nil))
	rlog := log.WithFields(log.Fields{"req": c.Request.RequestURI, "rid": rid, "ip": c.ClientIP()}).
		WithContext(context.WithValue(context.Background(), "time", time.Now()))
	rlog.WithTime(time.Now()).Debug("start")
	// -----------

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

	//trade, err := explorer.PoolService.FindSwapRoutePath(rlog, fromCoinId, toCoinId, tradeType, helpers.StringToBigInt(reqQuery.Amount))
	//if err != nil {
	//	errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
	//	return
	//}

	trade, err := explorer.PoolService.FindSwapRoutePathByNode(fromCoinId, toCoinId, reqQuery.TradeType, reqQuery.Amount)
	if err != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	rlog.WithTime(time.Now()).WithField("t", time.Since(rlog.Context.Value("time").(time.Time))).Debug("trade found")

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		coin, err := explorer.CoinRepository.FindByID(uint(t.CoinID))
		helpers.CheckErr(err)
		path[i] = coin
	}

	rlog.WithTime(time.Now()).WithField("t", time.Since(rlog.Context.Value("time").(time.Time))).Debug("result created")

	c.JSON(http.StatusOK, new(pool.RouteResource).Transform(path, trade))
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

	liquidityPools, err := explorer.PoolRepository.GetAll()
	helpers.CheckErr(err)

	var swapCoins []models.Coin

	wg := new(sync.WaitGroup)
	for _, pc := range poolCoins {
		if pc.ID == fromCoin.ID {
			continue
		}

		wg.Add(1)
		go func(pc models.Coin, wg *sync.WaitGroup) {
			defer wg.Done()
			if explorer.PoolService.IsSwapExists(liquidityPools, uint64(fromCoin.ID), uint64(pc.ID), reqQuery.Depth) {
				swapCoins = append(swapCoins, pc)
			}
		}(pc, wg)
	}
	wg.Wait()

	c.JSON(http.StatusOK, resource.TransformCollection(swapCoins, coins.IdResource{}))
}

func EstimateSwap(c *gin.Context) {
	explorer := c.MustGet("explorer").(*core.Explorer)

	// todo: remove, temp stats collector
	hasher := md5.New()
	hasher.Write([]byte(c.Request.RequestURI + time.Now().String() + c.ClientIP()))
	rid := hex.EncodeToString(hasher.Sum(nil))
	rlog := log.WithFields(log.Fields{"req": c.Request.RequestURI, "rid": rid, "ip": c.ClientIP()}).
		WithContext(context.WithValue(context.Background(), "time", time.Now()))
	rlog.WithTime(time.Now()).Debug("start")
	// -----------

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
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	// define trade type
	tradeType := swap.TradeTypeExactInput
	if reqQuery.TradeType == "output" {
		tradeType = swap.TradeTypeExactOutput
	}

	bancorAmount, bancorErr := explorer.SwapService.EstimateByBancor(coinFrom, coinTo, reqQuery.GetAmount(), tradeType)
	//trade, poolErr := explorer.PoolService.FindSwapRoutePath(rlog, uint64(coinFrom.ID), uint64(coinTo.ID), tradeType, reqQuery.GetAmount())
	trade, poolErr := explorer.PoolService.FindSwapRoutePathByNode(uint64(coinFrom.ID), uint64(coinTo.ID), reqQuery.TradeType, reqQuery.Amount)

	if poolErr != nil && bancorErr != nil {
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Route path not exists.", c)
		return
	}

	if bancorErr == nil && poolErr != nil {
		c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
		return
	}

	if bancorErr == nil && poolErr == nil {
		if tradeType == swap.TradeTypeExactInput && bancorAmount.Cmp(trade.OutputAmount.GetAmount()) >= 1 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}

		if tradeType == swap.TradeTypeExactOutput && bancorAmount.Cmp(trade.InputAmount.GetAmount()) <= 0 {
			c.JSON(http.StatusOK, new(pool.BancorResource).Transform(reqQuery.GetAmount(), bancorAmount, tradeType))
			return
		}
	}

	path := make([]models.Coin, len(trade.Route.Path))
	for i, t := range trade.Route.Path {
		path[i], _ = explorer.CoinRepository.FindByID(uint(t.CoinID))
	}

	c.JSON(http.StatusOK, new(pool.RouteResource).Transform(path, trade))
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
		errors.SetErrorResponse(http.StatusNotFound, http.StatusNotFound, "Pool not found.", c)
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

	pools, err := explorer.PoolRepository.GetAll()
	if err != nil {
		log.Panicf("failed to load pools: %s", err)
	}

	type cmcResource struct {
		BaseId      uint64 `json:"base_id"`
		BaseName    string `json:"base_name"`
		BaseSymbol  string `json:"base_symbol"`
		QuoteId     uint64 `json:"quote_id"`
		QuoteName   string `json:"quote_name"`
		QuoteSymbol string `json:"quote_symbol"`
		LastPrice   string `json:"last_price"`
		BaseVolume  string `json:"base_volume"`
		QuoteVolume string `json:"quote_volume"`
	}

	resources := make(map[string]cmcResource, len(pools))

	for _, p := range pools {
		ticker := fmt.Sprintf("%d_%d", p.FirstCoinId, p.SecondCoinId)

		price := new(big.Float).Quo(helpers.StrToBigFloat(p.FirstCoinVolume), helpers.StrToBigFloat(p.SecondCoinVolume))

		resources[ticker] = cmcResource{
			BaseId:      p.FirstCoinId,
			BaseName:    p.FirstCoin.Name,
			BaseSymbol:  p.FirstCoin.GetSymbol(),
			QuoteId:     p.SecondCoinId,
			QuoteName:   p.SecondCoin.Name,
			QuoteSymbol: p.SecondCoin.GetSymbol(),
			LastPrice:   helpers.Bip2Str(price),
			BaseVolume:  helpers.PipStr2Bip(p.FirstCoinVolume),
			QuoteVolume: helpers.PipStr2Bip(p.SecondCoinVolume),
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
