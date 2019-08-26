package balance

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-tools/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Service struct {
	baseCoin string
}

func NewService(baseCoin string) *Service {
	return &Service{baseCoin}
}

func (s *Service) GetBalanceSumByAddress(address *models.Address) *big.Int {
	sum := big.NewInt(0)

	if address == nil {
		return sum
	}

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

	return sum
}


