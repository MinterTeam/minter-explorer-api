package balance

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools/market"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
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

func (s *Service) GetTotalBalance(address *models.Address) *big.Float {
	sum := big.NewInt(0)
	for _, balance := range address.Balances {
		// just add base coin to sum
		if balance.Coin.Symbol == s.baseCoin {
			sum = sum.Add(sum, helpers.StringToBigInt(balance.Value))
			continue
		}

		// calculate the sale return value for custom coin
		sum = sum.Add(sum, formula.CalculateSaleReturn(
			helpers.StringToBigInt(balance.Coin.Volume),
			helpers.StringToBigInt(balance.Coin.ReserveBalance),
			uint(balance.Coin.Crr),
			helpers.StringToBigInt(balance.Value),
		))
	}

	return new(big.Float).SetInt(sum)
}

func (s *Service) GetTotalBalanceInUSD(sumInBasecoin *big.Float) *big.Float {
	return new(big.Float).Mul(sumInBasecoin, big.NewFloat(s.marketService.PriceChange.Price))
}
