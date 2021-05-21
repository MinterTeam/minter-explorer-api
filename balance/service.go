package balance

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/services"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/market"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Service struct {
	baseCoin      string
	marketService *market.Service
	swapService   *services.SwapService
}

func NewService(baseCoin string, marketService *market.Service, swapService *services.SwapService) *Service {
	return &Service{baseCoin, marketService, swapService}
}

func (s *Service) GetTotalBalance(address *models.Address) *big.Int {
	sum := big.NewInt(0)
	for _, balance := range address.Balances {
		sum = sum.Add(sum, s.swapService.EstimateInBip(balance.Coin, helpers.StringToBigInt(balance.Value)))
	}

	return sum
}

func (s *Service) GetStakeBalance(stakes []models.Stake) *big.Int {
	sum := big.NewInt(0)

	for _, stake := range stakes {
		if stake.IsKicked || stake.Coin.Crr == 0 {
			continue
		}

		// just add base coin to sum
		if stake.Coin.Symbol == s.baseCoin {
			sum = sum.Add(sum, helpers.StringToBigInt(stake.Value))
			continue
		}

		// calculate the sale return value for custom coin
		sum = sum.Add(sum, formula.CalculateSaleReturn(
			helpers.StringToBigInt(stake.Coin.Volume),
			helpers.StringToBigInt(stake.Coin.Reserve),
			uint(stake.Coin.Crr),
			helpers.StringToBigInt(stake.Value),
		))
	}

	return sum
}

func (s *Service) GetTotalBalanceInUSD(sumInBasecoin *big.Int) *big.Float {
	return new(big.Float).Mul(new(big.Float).SetInt(sumInBasecoin), big.NewFloat(s.marketService.PriceChange.Price))
}
