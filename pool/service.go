package pool

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetPoolLiquidityInBip(pool models.LiquidityPool) *big.Int {
	var (
		firstCoinInBip  = helpers.StringToBigInt(pool.FirstCoinVolume)
		secondCoinInBip = helpers.StringToBigInt(pool.SecondCoinVolume)
	)

	if pool.FirstCoinId != 0 && pool.FirstCoin.Crr != 0 {
		firstCoinInBip = formula.CalculateSaleReturn(
			helpers.StringToBigInt(pool.FirstCoin.Volume),
			helpers.StringToBigInt(pool.FirstCoin.Reserve),
			pool.FirstCoin.Crr,
			helpers.StringToBigInt(pool.FirstCoinVolume),
		)
	}

	if pool.SecondCoinId != 0 && pool.SecondCoin.Crr != 0 {
		secondCoinInBip = formula.CalculateSaleReturn(
			helpers.StringToBigInt(pool.SecondCoin.Volume),
			helpers.StringToBigInt(pool.SecondCoin.Reserve),
			pool.SecondCoin.Crr,
			helpers.StringToBigInt(pool.SecondCoinVolume),
		)
	}

	return new(big.Int).Add(firstCoinInBip, secondCoinInBip)
}
