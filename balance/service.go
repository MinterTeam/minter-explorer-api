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

type CoinBalance struct {
	Value *big.Int
	Coin  *models.Coin
}

func (s *Service) GetSumByBalancesAndStakes(addressBalances []*models.Balance, addressStakes []*models.Stake) *big.Float {
	var balances []CoinBalance

	// add to balances slice from address balances
	for _, balance := range addressBalances {
		balances = append(balances, CoinBalance{helpers.StringToBigInt(balance.Value), balance.Coin})
	}

	// add or sum balances from address stakes
	for _, stake := range addressStakes {
		coin := stake.Coin
		value := helpers.StringToBigInt(stake.Value)

		for key, balance := range balances {
			if balance.Coin.ID == stake.Coin.ID {
				value = new(big.Int).Add(balance.Value, value)
				balances = append(balances[:key], balances[key+1:]...)
				break
			}
		}

		balances = append(balances, CoinBalance{value, coin})
	}

	// sum overall balances in base coin
	sum := big.NewInt(0)
	for _, balance := range balances {
		// just add base coin to sum
		if balance.Coin.Symbol == s.baseCoin {
			sum = sum.Add(sum, balance.Value)
			continue
		}

		// calculate the sale return value for custom coin
		sum = sum.Add(sum, formula.CalculateSaleReturn(
			helpers.StringToBigInt(balance.Coin.Volume),
			helpers.StringToBigInt(balance.Coin.ReserveBalance),
			uint(balance.Coin.Crr),
			balance.Value,
		))
	}

	return new(big.Float).SetInt(sum)
}

func (s *Service) GetBalanceSumInUSDByBaseCoin(sumInBasecoin *big.Float) *big.Float {
	return new(big.Float).Mul(sumInBasecoin, big.NewFloat(s.marketService.PriceChange.Price))
}
