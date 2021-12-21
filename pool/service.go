package pool

import (
	"fmt"
	"github.com/MinterTeam/explorer-sdk/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	. "github.com/MinterTeam/minter-explorer-api/v2/errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sync"
	"time"
)

const (
	defaultTradeVolumeScale = "day"
)

type Service struct {
	repository            *Repository
	swap                  *swap.Service
	poolsLiquidity        map[uint64]*big.Float
	coinPrices            map[uint64]*big.Float
	tradeVolumes          map[uint64]TradeVolumes
	poolTokenToPoolMap    *sync.Map
	swapRoutes            *sync.Map
	pools                 []models.LiquidityPool
	poolsMap              *sync.Map
	coinToTrackedPoolsMap *sync.Map

	plmx sync.RWMutex
	cpmx sync.RWMutex
	tvmx sync.RWMutex

	lastTradeVolumesUpdate *time.Time
	node                   *resty.Client

	coinService *coins.Service
}

func NewService(repository *Repository, coinService *coins.Service) *Service {
	pools, _ := repository.GetAll()

	s := &Service{
		swap:                   swap.NewService(repository.db),
		coinService:            coinService,
		repository:             repository,
		pools:                  pools,
		poolsLiquidity:         make(map[uint64]*big.Float),
		coinPrices:             make(map[uint64]*big.Float),
		tradeVolumes:           make(map[uint64]TradeVolumes),
		swapRoutes:             new(sync.Map),
		poolsMap:               new(sync.Map),
		coinToTrackedPoolsMap:  new(sync.Map),
		poolTokenToPoolMap:     new(sync.Map),
		lastTradeVolumesUpdate: nil,
		node:                   resty.New(),
	}

	s.SavePoolsToMap()

	return s
}

func (s *Service) GetCoinPrice(coinId uint64) *big.Float {
	s.cpmx.RLock()
	defer s.cpmx.RUnlock()

	if amount, ok := s.coinPrices[coinId]; ok {
		return amount
	}

	return big.NewFloat(0)
}

func (s *Service) GetCoinPriceInBip(coinId uint64) *big.Float {
	s.cpmx.RLock()
	defer s.cpmx.RUnlock()

	if amount, ok := s.coinPrices[coinId]; ok {
		return new(big.Float).Quo(amount, s.coinPrices[0])
	}

	return big.NewFloat(0)
}

func (s *Service) IsSwapExists(pools []models.LiquidityPool, fromCoinId, toCoinId uint64, depth int) bool {
	_, err := s.swap.FindSwapRoutePathsByGraph(pools, fromCoinId, toCoinId, depth, 20)
	return err == nil
}

func (s *Service) OnNewBlock(block blocks.Resource) {
	if block.NumTxs == 0 {
		return
	}

	s.RunPoolUpdater()
}

func (s *Service) RunPoolUpdater() {
	defer Recovery()

	pools, err := s.repository.GetAll()
	if err != nil {
		log.Error(err)
	}

	s.pools = pools
	s.SavePoolsToMap()
	s.RunCoinPriceCalculation(pools)
	s.RunPoolTradeVolumesUpdater(pools)
}

func (s *Service) SavePoolsToMap() {
	for _, p := range s.pools {
		key := fmt.Sprintf("%d-%d", p.FirstCoinId, p.SecondCoinId)
		s.poolsMap.Store(key, p)
		s.poolTokenToPoolMap.Store(p.TokenId, p)
	}
}

func (s *Service) GetPoolByToken(token *models.Coin) *models.LiquidityPool {
	if lp, ok := s.poolTokenToPoolMap.Load(uint64(token.ID)); ok {
		pool := lp.(models.LiquidityPool)
		return &pool
	}

	return nil
}

func (s *Service) RunLiquidityCalculation(pools []models.LiquidityPool) {
	for _, p := range pools {
		liquidity := s.swap.GetPoolLiquidity(pools, p)
		s.plmx.Lock()
		s.poolsLiquidity[p.Id] = liquidity
		s.plmx.Unlock()
	}
}

func (s *Service) RunCoinPriceCalculation(pools []models.LiquidityPool) {
	verifiedCoins := s.coinService.GetVerifiedCoins()
	verifiedCoinIds := s.mapToVerifiedCoinIds(verifiedCoins)
	s.fillCoinToTrackedPoolsMap(pools, verifiedCoinIds)
	coinprices := s.computeCoinPrices(pools, s.computeTrackedCoinPrices(verifiedCoins))

	coinPricesMap := make(map[uint64]*big.Float)
	coinprices.Range(func(key, value interface{}) bool {
		coinPricesMap[uint64(key.(uint))] = value.(*big.Float)
		return true
	})

	s.cpmx.Lock()
	s.coinPrices = coinPricesMap
	s.cpmx.Unlock()
}

func (s *Service) mapToVerifiedCoinIds(verified []models.Coin) []uint64 {
	ids := []uint64{config.MUSD}
	for _, vc := range verified {
		ids = append(ids, uint64(vc.ID))
	}

	return ids
}

func (s *Service) fillCoinToTrackedPoolsMap(pools []models.LiquidityPool, verifiedCoinIds []uint64) {
	minLiquidity, _ := new(big.Int).SetString("10000000000000000000000", 10)
	coinToTrackedPoolsMap := make(map[uint64][]models.LiquidityPool)

	for _, p := range pools {
		if helpers.InArray(p.FirstCoinId, verifiedCoinIds) || helpers.InArray(p.SecondCoinId, verifiedCoinIds) {
			if helpers.StringToBigInt(p.LiquidityBip).Cmp(minLiquidity) == -1 {
				continue
			}

			if coinPools, ok := coinToTrackedPoolsMap[p.FirstCoinId]; ok {
				coinToTrackedPoolsMap[p.FirstCoinId] = append(coinPools, p)
			} else {
				coinToTrackedPoolsMap[p.FirstCoinId] = []models.LiquidityPool{p}
			}

			if coinPools, ok := coinToTrackedPoolsMap[p.SecondCoinId]; ok {
				coinToTrackedPoolsMap[p.SecondCoinId] = append(coinPools, p)
			} else {
				coinToTrackedPoolsMap[p.SecondCoinId] = []models.LiquidityPool{p}
			}
		}
	}

	for id, coinPools := range coinToTrackedPoolsMap {
		s.coinToTrackedPoolsMap.Store(id, coinPools)
	}
}

func (s *Service) computeTrackedCoinPrices(verifiedCoins []models.Coin) *sync.Map {
	coinPrices := new(sync.Map)
	coinPrices.Store(uint(config.MUSD), big.NewFloat(1))

	// move bip to start
	verifiedCoinsSorted := make([]models.Coin, len(verifiedCoins))
	for i, p := range verifiedCoins {
		if p.ID == 0 {
			verifiedCoinsSorted = append(append([]models.Coin{p}, verifiedCoins[:i]...), verifiedCoins[i+1:]...)
		}
	}

	for _, p := range verifiedCoinsSorted {
		pid := uint64(p.ID)
		if helpers.InArray(pid, config.StableCoinIds) {
			coinPrices.Store(p.ID, big.NewFloat(1))
			continue
		}

		coinPools, ok := s.coinToTrackedPoolsMap.Load(pid)
		if !ok {
			continue
		}

		for _, cp := range coinPools.([]models.LiquidityPool) {
			if cp.FirstCoinId == 0 && p.ID != 0 {
				firstCoinVolume := helpers.StrToBigFloat(cp.FirstCoinVolume)
				secondCoinVolume := helpers.StrToBigFloat(cp.SecondCoinVolume)
				bipPrice, _ := coinPrices.Load(uint(0))
				bipToCoin := new(big.Float).Quo(secondCoinVolume, firstCoinVolume)
				usdPrice := new(big.Float).Quo(bipPrice.(*big.Float), bipToCoin)
				coinPrices.Store(p.ID, usdPrice)
				break
			}

			if helpers.InArray(cp.FirstCoinId, config.StableCoinIds) || helpers.InArray(cp.SecondCoinId, config.StableCoinIds) {
				firstCoinVolume := helpers.StrToBigFloat(cp.FirstCoinVolume)
				secondCoinVolume := helpers.StrToBigFloat(cp.SecondCoinVolume)

				if cp.FirstCoinId == pid {
					coinPrices.Store(p.ID, new(big.Float).Quo(secondCoinVolume, firstCoinVolume))
					break
				}

				coinPrices.Store(p.ID, new(big.Float).Quo(firstCoinVolume, secondCoinVolume))
				break
			}
		}
	}

	return coinPrices
}

func (s *Service) computeCoinPrices(pools []models.LiquidityPool, coinPrices *sync.Map) *sync.Map {
	wg := new(sync.WaitGroup)
	for _, p := range pools {
		wg.Add(1)
		go func(p models.LiquidityPool, wg *sync.WaitGroup) {
			defer wg.Done()
			if _, ok := coinPrices.Load(uint(p.FirstCoinId)); !ok {
				if coinPools, ok := s.coinToTrackedPoolsMap.Load(p.FirstCoinId); ok {
					cp := coinPools.([]models.LiquidityPool)[0]
					firstCoinVolume := helpers.StrToBigFloat(cp.FirstCoinVolume)
					secondCoinVolume := helpers.StrToBigFloat(cp.SecondCoinVolume)

					if cp.FirstCoinId == p.FirstCoinId {
						if price, ok := coinPrices.Load(uint(cp.SecondCoinId)); ok {
							temp := new(big.Float).Quo(secondCoinVolume, firstCoinVolume)
							coinPrices.Store(uint(p.FirstCoinId), temp.Mul(temp, price.(*big.Float)))
						}
					} else {
						if price, ok := coinPrices.Load(uint(cp.FirstCoinId)); ok {
							temp := new(big.Float).Quo(firstCoinVolume, secondCoinVolume)
							coinPrices.Store(uint(p.FirstCoinId), temp.Mul(temp, price.(*big.Float)))
						}
					}
				}
			}

			if _, ok := coinPrices.Load(uint(p.SecondCoinId)); !ok {
				if coinPools, ok := s.coinToTrackedPoolsMap.Load(p.SecondCoinId); ok {
					cp := coinPools.([]models.LiquidityPool)[0]
					firstCoinVolume := helpers.StrToBigFloat(cp.FirstCoinVolume)
					secondCoinVolume := helpers.StrToBigFloat(cp.SecondCoinVolume)

					if cp.FirstCoinId == p.SecondCoinId {
						if price, ok := coinPrices.Load(uint(cp.SecondCoinId)); ok {
							temp := new(big.Float).Quo(secondCoinVolume, firstCoinVolume)
							coinPrices.Store(uint(p.SecondCoinId), temp.Mul(temp, price.(*big.Float)))
						}
					} else {
						if price, ok := coinPrices.Load(uint(cp.FirstCoinId)); ok {
							temp := new(big.Float).Quo(firstCoinVolume, secondCoinVolume)
							coinPrices.Store(uint(p.SecondCoinId), temp.Mul(temp, price.(*big.Float)))
						}
					}
				}
			}
		}(p, wg)
	}
	wg.Wait()

	return coinPrices
}

func (s *Service) RunPoolTradeVolumesUpdater(pools []models.LiquidityPool) {
	t := time.Now()
	if s.lastTradeVolumesUpdate != nil && s.lastTradeVolumesUpdate.Add(10*time.Minute).Before(t) {
		return
	}

	volumes1d, err := s.getTradeVolumesLastNDays(pools, 1)
	if err != nil {
		log.Errorf("failed to update last day trade volumes: %s", err)
		return
	}

	volumes30d, err := s.getTradeVolumesLastNDays(pools, 30)
	if err != nil {
		log.Errorf("failed to update last month trade volumes: %s", err)
		return
	}

	for pId, v := range volumes1d {
		s.tvmx.Lock()
		s.tradeVolumes[pId] = TradeVolumes{Day: v, Month: volumes30d[pId]}
		s.tvmx.Unlock()
	}

	s.lastTradeVolumesUpdate = &t
}

func (s *Service) GetTradesVolume(pool models.LiquidityPool, scale *string, startTime *time.Time) ([]TradeVolume, error) {
	dateScale := defaultTradeVolumeScale
	if scale != nil {
		dateScale = *scale
	}

	tradesVolume, err := s.repository.GetPoolTradesVolume(pool, dateScale, startTime)
	if err != nil {
		return nil, err
	}

	bipPrice := s.GetCoinPriceInBip(pool.FirstCoinId)
	trades := make([]TradeVolume, len(tradesVolume))

	for i, tv := range tradesVolume {
		bipVolume := helpers.Pip2Bip(helpers.StringToBigInt(tv.FirstCoinVolume))
		if pool.FirstCoinId != 0 {
			bipVolume = getVolumeInBip(bipPrice, tv.FirstCoinVolume)
		}

		trades[i] = TradeVolume{
			Date:             tv.Date,
			FirstCoinVolume:  tv.FirstCoinVolume,
			SecondCoinVolume: tv.SecondCoinVolume,
			BipVolume:        bipVolume,
		}
	}

	return trades, nil
}

func (s *Service) getLastMonthTradesVolume(pool models.LiquidityPool) (*TradeVolume, error) {
	startTime := time.Now().AddDate(0, -1, 0)
	tv, err := s.repository.GetPoolTradeVolumeByTimeRange(pool, startTime)
	if err != nil {
		return nil, err
	}

	if tv == nil {
		return &TradeVolume{
			BipVolume: big.NewFloat(0),
		}, nil
	}

	bipPrice := s.GetCoinPriceInBip(pool.FirstCoinId)
	bipVolume := helpers.Pip2Bip(helpers.StringToBigInt(tv.FirstCoinVolume))
	if pool.FirstCoinId != 0 {
		bipVolume = getVolumeInBip(bipPrice, tv.FirstCoinVolume)
	}

	return &TradeVolume{
		FirstCoinVolume:  tv.FirstCoinVolume,
		SecondCoinVolume: tv.SecondCoinVolume,
		BipVolume:        bipVolume,
	}, nil
}

func (s *Service) getTradeVolumesLastNDays(pools []models.LiquidityPool, days int) (map[uint64]TradeVolume, error) {
	poolsMap := make(map[uint64]models.LiquidityPool, len(pools))
	for _, p := range pools {
		poolsMap[p.Id] = p
	}

	startTime := time.Now().AddDate(0, 0, -days)
	volumes, err := s.repository.GetPoolsTradeVolumeByTimeRange(pools, startTime)
	if err != nil {
		return nil, err
	}

	tvs := make(map[uint64]TradeVolume, len(pools))
	for _, v := range volumes {
		bipPrice := s.GetCoinPriceInBip(poolsMap[v.PoolId].FirstCoinId)
		bipVolume := helpers.Pip2Bip(helpers.StringToBigInt(v.FirstCoinVolume))
		if poolsMap[v.PoolId].FirstCoinId != 0 {
			bipVolume = getVolumeInBip(bipPrice, v.FirstCoinVolume)
		}

		tvs[v.PoolId] = TradeVolume{
			FirstCoinVolume:  v.FirstCoinVolume,
			SecondCoinVolume: v.SecondCoinVolume,
			BipVolume:        bipVolume,
		}
	}

	for _, p := range pools {
		if _, ok := tvs[p.Id]; !ok {
			tvs[p.Id] = TradeVolume{
				BipVolume: big.NewFloat(0),
			}
		}
	}

	return tvs, nil
}

func (s *Service) GetLastMonthTradesVolume(pool models.LiquidityPool) *TradeVolume {
	s.tvmx.RLock()
	defer s.tvmx.RUnlock()

	if tv, ok := s.tradeVolumes[pool.Id]; ok {
		return &tv.Month
	}

	return &TradeVolume{
		FirstCoinVolume:  "0",
		SecondCoinVolume: "0",
		BipVolume:        big.NewFloat(0),
	}
}

func (s *Service) GetLastDayTradesVolume(pool models.LiquidityPool) *TradeVolume {
	s.tvmx.RLock()
	defer s.tvmx.RUnlock()

	if tv, ok := s.tradeVolumes[pool.Id]; ok {
		return &tv.Day
	}

	return &TradeVolume{
		FirstCoinVolume:  "0",
		SecondCoinVolume: "0",
		BipVolume:        big.NewFloat(0),
	}
}

func (s *Service) GetPools() []models.LiquidityPool {
	return s.pools
}
