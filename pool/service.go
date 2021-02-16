package pool

import (
	"errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/swap"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"github.com/starwander/goraph"
	"math/big"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository}
}

func (s *Service) GetPoolLiquidityInBip(pool models.LiquidityPool) *big.Int {
	var (
		firstCoinInBip  = helpers.StringToBigInt(pool.FirstCoinVolume)
		secondCoinInBip = helpers.StringToBigInt(pool.SecondCoinVolume)
	)

	if pool.FirstCoin.Crr == 0 || pool.SecondCoin.Crr == 0 {
		return big.NewInt(0)
	}

	if pool.FirstCoinId != 0 {
		firstCoinInBip = formula.CalculateSaleReturn(
			helpers.StringToBigInt(pool.FirstCoin.Volume),
			helpers.StringToBigInt(pool.FirstCoin.Reserve),
			pool.FirstCoin.Crr,
			helpers.StringToBigInt(pool.FirstCoinVolume),
		)
	}

	if pool.SecondCoinId != 0 {
		secondCoinInBip = formula.CalculateSaleReturn(
			helpers.StringToBigInt(pool.SecondCoin.Volume),
			helpers.StringToBigInt(pool.SecondCoin.Reserve),
			pool.SecondCoin.Crr,
			helpers.StringToBigInt(pool.SecondCoinVolume),
		)
	}

	return new(big.Int).Add(firstCoinInBip, secondCoinInBip)
}

func (s *Service) FindSwapRoutePath(fromCoinId, toCoinId uint64, tradeType swap.TradeType, amount *big.Int) (*swap.Trade, error) {
	liquidityPools, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	paths, err := s.FindSwapRoutePathsByGraph(liquidityPools, fromCoinId, toCoinId)
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
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 3},
		)
	} else {
		trades, err = swap.GetBestTradeExactOut(
			pairs,
			swap.NewToken(fromCoinId),
			swap.NewTokenAmount(swap.NewToken(toCoinId), amount),
			swap.TradeOptions{MaxNumResults: 1, MaxHops: 3},
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

func (s *Service) FindSwapRoutePathsByGraph(pools []models.LiquidityPool, fromCoinId, toCoinId uint64) ([][]goraph.ID, error) {
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

	return paths, nil
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
