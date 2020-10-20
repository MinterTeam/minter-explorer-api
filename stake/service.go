package stake

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) PrepareStakesModels(stakes []models.Stake) ([]models.Stake, error) {
	for i, stake := range stakes {
		if stake.IsKicked {
			bipValue, err := s.calculateBipValue(stake, stake.Coin)
			if err != nil {
				return nil, err
			}

			stakes[i].BipValue = bipValue.String()
		}
	}

	return stakes, nil
}

func (s *Service) calculateBipValue(stake models.Stake, coin *models.Coin) (*big.Int, error) {
	if coin.ID == 0 {
		return helpers.StringToBigInt(stake.Value), nil
	}

	totalStakeStr, err := s.repository.GetSumValueByCoin(coin.ID)
	if err != nil {
		return nil, err
	}

	totalStake := helpers.StringToBigInt(totalStakeStr)
	stakeValue := helpers.StringToBigInt(stake.Value)
	coinVolume := helpers.StringToBigInt(coin.Volume)
	coinReserve := helpers.StringToBigInt(coin.Reserve)

	freeFloat := new(big.Int).Sub(coinVolume, totalStake)
	freeFloatInBip := formula.CalculateSaleReturn(coinVolume, coinReserve, coin.Crr, freeFloat)
	delegatedBipValue := new(big.Int).Sub(coinReserve, freeFloatInBip)

	return new(big.Int).Div(new(big.Int).Mul(stakeValue, delegatedBipValue), totalStake), nil
}
