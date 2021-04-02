package pool

import (
	"errors"
	"github.com/MinterTeam/minter-explorer-api/v2/blocks"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/swap"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	log "github.com/sirupsen/logrus"
	"github.com/starwander/goraph"
	"math/big"
)

const (
	defaultTradeVolumeScale = "day"
)

type Service struct {
	repository     *Repository
	poolsLiquidity map[uint64]*big.Int
	coinPrices     map[uint64]*big.Float
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository:     repository,
		poolsLiquidity: make(map[uint64]*big.Int),
		coinPrices:     make(map[uint64]*big.Float),
	}
}

func (s *Service) GetPoolLiquidityInBip(pool models.LiquidityPool) *big.Int {
	if liquidity, ok := s.poolsLiquidity[pool.Id]; ok {
		return liquidity
	}

	return big.NewInt(0)
}

func (s *Service) GetCoinPriceInBip(coinId uint64) *big.Float {
	if amount, ok := s.coinPrices[coinId]; ok {
		return amount
	}

	return big.NewFloat(0)
}

func (s *Service) FindSwapRoutePath(fromCoinId, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	pools, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	return s.FindSwapRoutePathByPools(pools, fromCoinId, toCoinId, tradeType, amount)
}

func (s *Service) FindSwapRoutePathByPools(liquidityPools []models.LiquidityPool, fromCoinId, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	paths, err := s.FindSwapRoutePathsByGraph(liquidityPools, fromCoinId, toCoinId, 1000)
	if err != nil {
		return nil, err
	}

	pools, err := s.getPathsRelatedPools(paths)
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

func (s *Service) FindSwapRoutePathsByGraph(pools []models.LiquidityPool, fromCoinId, toCoinId uint64, depth int) ([][]goraph.ID, error) {
	graph := goraph.NewGraph()
	for _, pool := range pools {
		graph.AddVertex(pool.FirstCoinId, pool.FirstCoin)
		graph.AddVertex(pool.SecondCoinId, pool.SecondCoin)
		graph.AddEdge(pool.FirstCoinId, pool.SecondCoinId, 1, nil)
		graph.AddEdge(pool.SecondCoinId, pool.FirstCoinId, 1, nil)
	}

	_, paths, err := graph.Yen(fromCoinId, toCoinId, 1000)
	if err != nil {
		return nil, err
	}

	if len(paths[0]) == 0 {
		return nil, errors.New("path not found")
	}

	if depth == 0 {
		return paths, nil
	}

	var result [][]goraph.ID
	for _, path := range paths {
		if len(path) > depth+1 || len(path) == 0 {
			break
		}

		result = append(result, path)
	}

	if len(result) == 0 {
		return nil, errors.New("path not found")
	}

	return result, nil
}

func (s *Service) getPathsRelatedPools(paths [][]goraph.ID) ([]models.LiquidityPool, error) {
	pools := make([]models.LiquidityPool, 0)
	for _, path := range paths {
		if len(path) == 0 {
			break
		}

		for i, coinID := range path {
			if i == 0 {
				continue
			}

			p, err := s.repository.Find(path[i-1].(uint64), coinID.(uint64))
			if err != nil {
				return nil, err
			}

			isExists := false
			for _, p2 := range pools {
				if p.Id == p2.Id {
					isExists = true
				}
			}

			if !isExists {
				pools = append(pools, p)
			}
		}
	}

	return pools, nil
}

// TODO: refactor and optimize
func (s *Service) OnNewBlock(block blocks.Resource) {
	if block.NumTxs == 0 {
		return
	}

	pools, err := s.repository.GetAll()
	if err != nil {
		log.Error(err)
	}

	for _, p := range pools {
		firstCoinVolume := helpers.StringToBigInt(p.FirstCoinVolume)
		secondCoinVolume := helpers.StringToBigInt(p.SecondCoinVolume)

		if p.FirstCoinId == 0 {
			s.poolsLiquidity[p.Id] = new(big.Int).Mul(firstCoinVolume, big.NewInt(2))
			continue
		}

		trade, err := s.FindSwapRoutePathByPools(pools, p.FirstCoinId, uint64(0), swap.TradeTypeExactInput, firstCoinVolume)
		if err != nil || trade == nil {
			trade, err = s.FindSwapRoutePathByPools(pools, p.SecondCoinId, uint64(0), swap.TradeTypeExactInput, secondCoinVolume)
			if err != nil || trade == nil {
				s.poolsLiquidity[p.Id] = big.NewInt(0)
				continue
			}

			price := swap.NewPrice(trade.InputAmount.Token, trade.OutputAmount.Token, trade.InputAmount.GetAmount(), trade.OutputAmount.GetAmount()).Value
			s.coinPrices[p.SecondCoinId] = helpers.Pip2Bip(price)
			s.coinPrices[p.FirstCoinId] = new(big.Float).Mul(new(big.Float).Quo(
				new(big.Float).SetInt(secondCoinVolume),
				new(big.Float).SetInt(firstCoinVolume),
			), helpers.Pip2Bip(price))
		} else {
			price := swap.NewPrice(trade.InputAmount.Token, trade.OutputAmount.Token, trade.InputAmount.GetAmount(), trade.OutputAmount.GetAmount()).Value
			s.coinPrices[p.FirstCoinId] = helpers.Pip2Bip(price)
			s.coinPrices[p.SecondCoinId] = new(big.Float).Mul(new(big.Float).Quo(
				new(big.Float).SetInt(firstCoinVolume),
				new(big.Float).SetInt(secondCoinVolume),
			), helpers.Pip2Bip(price))
		}

		s.poolsLiquidity[p.Id] = new(big.Int).Mul(trade.OutputAmount.GetAmount(), big.NewInt(2))
	}
}

func (s *Service) GetTradesVolume(pool models.LiquidityPool, scale *string) ([]TradeVolume, error) {
	dateScale := defaultTradeVolumeScale
	if scale != nil {
		dateScale = *scale
	}

	tradesVolume, err := s.repository.GetPoolTradesVolume(pool, dateScale)
	if err != nil {
		return nil, err
	}

	bipPrice := s.GetCoinPriceInBip(pool.FirstCoinId)
	trades := make([]TradeVolume, len(tradesVolume))

	for i, tv := range tradesVolume {
		firstCoinBaseVolume := helpers.Pip2Bip(helpers.StringToBigInt(tv.FirstCoinVolume))
		bipVolume := new(big.Float).Mul(firstCoinBaseVolume, bipPrice)

		trades[i] = TradeVolume{
			Date:             tv.Date,
			FirstCoinVolume:  tv.FirstCoinVolume,
			SecondCoinVolume: tv.SecondCoinVolume,
			BipVolume:        bipVolume,
		}
	}

	return trades, nil
}
