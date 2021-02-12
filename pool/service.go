package pool

import (
	"errors"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
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

type myVertex struct {
	id     uint64
	outTo  map[uint64]float64
	inFrom map[uint64]float64
}

type myEdge struct {
	from   uint64
	to     uint64
	weight float64
}

func (vertex *myVertex) ID() goraph.ID {
	return vertex.id
}

func (vertex *myVertex) Edges() (edges []goraph.Edge) {
	edges = make([]goraph.Edge, len(vertex.outTo)+len(vertex.inFrom))
	i := 0
	for to, weight := range vertex.outTo {
		edges[i] = &myEdge{vertex.id, to, weight}
		i++
	}
	for from, weight := range vertex.inFrom {
		edges[i] = &myEdge{from, vertex.id, weight}
		i++
	}
	return
}

func (edge *myEdge) Get() (goraph.ID, goraph.ID, float64) {
	return edge.from, edge.to, edge.weight
}

func (s *Service) FindSwapRoutePath(fromCoinId, toCoinId uint64) ([]*models.Coin, error) {
	pools, err := s.repository.GetAll()
	if err != nil {
		return nil, err
	}

	graph := goraph.NewGraph()
	for _, pool := range pools {
		graph.AddVertex(pool.FirstCoinId, pool.FirstCoin)
		graph.AddVertex(pool.SecondCoinId, pool.SecondCoin)
		graph.AddEdge(pool.FirstCoinId, pool.SecondCoinId, 1, nil)
		graph.AddEdge(pool.SecondCoinId, pool.FirstCoinId, 1, nil)
	}

	_, path, err := graph.Yen(fromCoinId, toCoinId, 1)
	if err != nil {
		return nil, err
	}

	bestPath := path[0]
	if len(bestPath) == 0 {
		return nil, errors.New("path not found")
	}

	coinPath := make([]*models.Coin, len(bestPath))
	for i, id := range bestPath {
		v, _ := graph.GetVertex(id)
		coinPath[i] = v.(*models.Coin)
	}

	return coinPath, nil
}
