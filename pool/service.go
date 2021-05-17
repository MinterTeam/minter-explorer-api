package pool

import (
	"errors"
	"github.com/MinterTeam/explorer-sdk/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	log "github.com/sirupsen/logrus"
	"github.com/starwander/goraph"
	"math/big"
	"sync"
	"time"
)

const (
	defaultTradeVolumeScale = "day"
)

type Service struct {
	repository     *Repository
	swap           *swap.Service
	poolsLiquidity map[uint64]*big.Float
	coinPrices     map[uint64]*big.Float

	plmx sync.RWMutex
	cpmx sync.RWMutex
}

func NewService(repository *Repository) *Service {
	return &Service{
		swap:           swap.NewService(repository.db),
		repository:     repository,
		poolsLiquidity: make(map[uint64]*big.Float),
		coinPrices:     make(map[uint64]*big.Float),
	}
}

func (s *Service) GetPoolLiquidityInBip(pool models.LiquidityPool) *big.Float {
	s.plmx.RLock()
	defer s.plmx.RUnlock()

	if liquidity, ok := s.poolsLiquidity[pool.Id]; ok {
		return liquidity
	}

	return big.NewFloat(0)
}

func (s *Service) GetCoinPriceInBip(coinId uint64) *big.Float {
	s.cpmx.RLock()
	defer s.cpmx.RUnlock()

	if amount, ok := s.coinPrices[coinId]; ok {
		return amount
	}

	return big.NewFloat(0)
}

func (s *Service) IsSwapExists(pools []models.LiquidityPool, fromCoinId, toCoinId uint64, depth int) bool {
	_, err := s.swap.FindSwapRoutePathsByGraph(pools, fromCoinId, toCoinId, depth)
	return err == nil
}

func (s *Service) FindSwapRoutePath(fromCoinId, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	pools, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return s.FindSwapRoutePathByPools(pools, fromCoinId, toCoinId, tradeType, amount)
}

func (s *Service) FindSwapRoutePathByPools(liquidityPools []models.LiquidityPool, fromCoinId, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	paths, err := s.swap.FindSwapRoutePathsByGraph(liquidityPools, fromCoinId, toCoinId, 5)
	if err != nil {
		return nil, err
	}

	pools, err := s.getPathsRelatedPools(liquidityPools, paths)
	if err != nil {
		return nil, err
	}

	pairs := make([]swap.Pair, 0)
	for _, p := range pools {
		pairs = append(pairs, swap.NewPair(
			swap.NewTokenAmount(swap.NewToken(p.FirstCoinId), helpers.StringToBigInt(p.FirstCoinVolume)),
			swap.NewTokenAmount(swap.NewToken(p.SecondCoinId), helpers.StringToBigInt(p.SecondCoinVolume)),
		))
	}

	var trades []swap.Trade
	if tradeType == swap.TradeTypeExactInput {
		trades, err = swap.GetBestTradeExactIn(
			pairs,
			swap.NewToken(toCoinId),
			swap.NewTokenAmount(swap.NewToken(fromCoinId), amount),
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 4},
		)
	} else {
		trades, err = swap.GetBestTradeExactOut(
			pairs,
			swap.NewToken(fromCoinId),
			swap.NewTokenAmount(swap.NewToken(toCoinId), amount),
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 4},
		)
	}

	if err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return nil, errors.New("path not found")
	}

	return &trades[0], nil
}

func (s *Service) getPathsRelatedPools(pools []models.LiquidityPool, paths [][]goraph.ID) ([]models.LiquidityPool, error) {
	related := make([]models.LiquidityPool, 0)
	for _, path := range paths {
		if len(path) == 0 {
			break
		}

		for i := range path {
			if i == 0 {
				continue
			}

			firstCoinId, secondCoinId := path[i-1].(uint64), path[i].(uint64)
			pchan := make(chan models.LiquidityPool)
			for _, lp := range pools {
				go func(lp models.LiquidityPool) {
					if (lp.FirstCoinId == firstCoinId && lp.SecondCoinId == secondCoinId) || (lp.FirstCoinId == secondCoinId && lp.SecondCoinId == firstCoinId) {
						pchan <- lp
					}
				}(lp)
			}

			p := <-pchan

			isExists := false
			for _, p2 := range related {
				if p.Id == p2.Id {
					isExists = true
				}
			}

			if !isExists {
				related = append(related, p)
			}
		}
	}

	return related, nil
}

func (s *Service) OnNewBlock(block blocks.Resource) {
	if block.NumTxs == 0 {
		return
	}

	s.RunPoolUpdater()
}

func (s *Service) RunPoolUpdater() {
	pools, err := s.repository.GetAll()
	if err != nil {
		log.Error(err)
	}

	s.RunLiquidityCalculation(pools)
	s.RunCoinPriceCalculation(pools)
}

func (s *Service) RunLiquidityCalculation(pools []models.LiquidityPool) {
	for _, p := range pools {
		//s.plmx.Lock()
		s.poolsLiquidity[p.Id] = s.swap.GetPoolLiquidity(pools, p)
		//s.plmx.Unlock()
	}
}

func (s *Service) RunCoinPriceCalculation(pools []models.LiquidityPool) {
	coinPrice := make(map[uint64]*big.Float)

	for _, p := range pools {
		firstCoinVolume := helpers.StrToBigFloat(p.FirstCoinVolume)
		secondCoinVolume := helpers.StrToBigFloat(p.SecondCoinVolume)

		if p.FirstCoinId == 0 {
			if firstCoinVolume.Cmp(big.NewFloat(1e20)) >= 0 {
				coinPrice[p.SecondCoinId] = new(big.Float).Quo(firstCoinVolume, secondCoinVolume)
			}
		}
	}

	for k, v := range coinPrice {
		//s.cpmx.Lock()
		s.coinPrices[k] = v
		//s.cpmx.Unlock()
	}
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

func (s *Service) GetLastMonthTradesVolume(pool models.LiquidityPool) (*TradeVolume, error) {
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

func (s *Service) GetPoolsLastMonthTradesVolume(pools []models.LiquidityPool) (map[uint64]TradeVolume, error) {
	poolsMap := make(map[uint64]models.LiquidityPool, len(pools))
	for _, p := range pools {
		poolsMap[p.Id] = p
	}

	startTime := time.Now().AddDate(0, -1, 0)
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
