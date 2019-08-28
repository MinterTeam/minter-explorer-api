package balance

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools/market"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Service struct {
	baseCoin      string
	marketService *market.Service
}

func NewService(baseCoin string, marketService *market.Service) *Service {
	return &Service{baseCoin, marketService}
}

func (s *Service) GetBalanceSumByAddress(address *models.Address) *big.Float {
	if address == nil {
		return big.NewFloat(0)
	}

	sum := big.NewInt(0)
	for _, balance := range address.Balances {
		if balance.Coin.Symbol == s.baseCoin {
			sum = sum.Add(sum, helpers.StringToBigInt(balance.Value))
			continue
		}

		sum = sum.Add(sum, formula.CalculateSaleReturn(
			helpers.StringToBigInt(balance.Coin.Volume),
			helpers.StringToBigInt(balance.Coin.ReserveBalance),
			uint(balance.Coin.Crr),
			helpers.StringToBigInt(balance.Value),
		))
	}

	return new(big.Float).SetInt(sum)
}

func (s *Service) GetBalanceSumInUSDByBaseCoin(sumInBasecoin *big.Float) *big.Float {
	return new(big.Float).Mul(sumInBasecoin, big.NewFloat(s.marketService.PriceChange.Price))
}
