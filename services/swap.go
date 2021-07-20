package services

import (
	"errors"
	"github.com/MinterTeam/explorer-sdk/swap"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type SwapService struct {
	poolService *pool.Service
}

// TODO: refactor
var Swap *SwapService

func NewSwapService(poolService *pool.Service) *SwapService {
	return &SwapService{poolService}
}

func (s *SwapService) EstimateInBip(coin *models.Coin, swapAmount *big.Int) *big.Int {
	bancorAmount := s.EstimateInBipByBancor(coin, swapAmount)
	poolAmount := s.estimateInBipByPools(coin, swapAmount)

	if bancorAmount.Cmp(poolAmount) >= 0 {
		return bancorAmount
	}

	return poolAmount
}

func (s *SwapService) EstimateInBipByBancor(coin *models.Coin, swapAmount *big.Int) *big.Int {
	if coin.Crr == 0 {
		return big.NewInt(0)
	}

	if coin.ID == 0 {
		return swapAmount
	}

	return formula.CalculateSaleReturn(
		helpers.StringToBigInt(coin.Volume),
		helpers.StringToBigInt(coin.Reserve),
		coin.Crr,
		swapAmount,
	)
}

func (s *SwapService) estimateInBipByPools(coin *models.Coin, swapAmount *big.Int) *big.Int {
	price := s.poolService.GetCoinPriceInBip(uint64(coin.ID))
	bip := new(big.Float).Mul(helpers.Pip2Bip(swapAmount), price)
	return helpers.Bip2Pip(bip)
}

func (s *SwapService) EstimateByBancor(coinFrom models.Coin, coinTo models.Coin, swapAmount *big.Int, tradeType swap.TradeType) (*big.Int, error) {
	if tradeType == swap.TradeTypeExactInput {
		return s.EstimateSellByBancor(coinFrom, coinTo, swapAmount)
	}

	return s.EstimateBuyByBancor(coinFrom, coinTo, swapAmount)
}

func (s *SwapService) EstimateSellByBancor(coinFrom models.Coin, coinTo models.Coin, swapAmount *big.Int) (*big.Int, error) {
	if coinFrom.Type != models.CoinTypeBase || coinTo.Type != models.CoinTypeBase {
		return nil, errors.New("coin has no reserve")
	}

	if coinFrom.ID != 0 {
		swapAmount = formula.CalculateSaleReturn(
			helpers.StringToBigInt(coinFrom.Volume),
			helpers.StringToBigInt(coinFrom.Reserve),
			coinFrom.Crr,
			swapAmount,
		)

		if !s.checkCoinReserveUnderflow(coinFrom, swapAmount) {
			return nil, errors.New("coin reserve underflow")
		}
	}

	if coinTo.ID != 0 {
		swapAmount = formula.CalculatePurchaseReturn(
			helpers.StringToBigInt(coinTo.Volume),
			helpers.StringToBigInt(coinTo.Reserve),
			coinTo.Crr,
			swapAmount,
		)

		if !s.checkCoinReserveUnderflow(coinTo, swapAmount) {
			return nil, errors.New("coin reserve underflow")
		}
	}

	return swapAmount, nil
}

func (s *SwapService) EstimateBuyByBancor(coinFrom models.Coin, coinTo models.Coin, swapAmount *big.Int) (*big.Int, error) {
	if coinFrom.Type != models.CoinTypeBase || coinTo.Type != models.CoinTypeBase {
		return nil, errors.New("coin has no reserve")
	}

	if coinTo.ID != 0 {
		swapAmount = formula.CalculatePurchaseAmount(
			helpers.StringToBigInt(coinTo.Volume),
			helpers.StringToBigInt(coinTo.Reserve),
			coinTo.Crr,
			swapAmount,
		)

		if !s.checkCoinReserveUnderflow(coinTo, swapAmount) {
			return nil, errors.New("coin reserve underflow")
		}
	}

	if coinFrom.ID != 0 {
		swapAmount = formula.CalculateSaleAmount(
			helpers.StringToBigInt(coinFrom.Volume),
			helpers.StringToBigInt(coinFrom.Reserve),
			coinFrom.Crr,
			swapAmount,
		)

		if !s.checkCoinReserveUnderflow(coinFrom, swapAmount) {
			return nil, errors.New("coin reserve underflow")
		}
	}

	return swapAmount, nil
}

func (s *SwapService) checkCoinReserveUnderflow(coin models.Coin, delta *big.Int) bool {
	total := big.NewInt(0).Sub(helpers.StringToBigInt(coin.Reserve), delta)

	if total.Cmp(helpers.Bip2Pip(big.NewFloat(10000))) == -1 {
		return false
	}

	return true
}