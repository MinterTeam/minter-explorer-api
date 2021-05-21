package address

import (
	"github.com/MinterTeam/minter-explorer-api/v2/balance"
	"github.com/MinterTeam/minter-explorer-api/v2/stake"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Service struct {
	rp              *Repository
	stakeRepository *stake.Repository
	balanceService  *balance.Service
}

func NewService(rp *Repository, stakeRepository *stake.Repository, balanceService *balance.Service) *Service {
	return &Service{
		rp:              rp,
		stakeRepository: stakeRepository,
		balanceService:  balanceService,
	}
}

func (s *Service) GetBalance(minterAddress string, isTotal bool) (*Balance, error) {
	model := s.rp.GetByAddress(minterAddress)

	if isTotal {
		return s.GetTotalBalance(model)
	}

	return &Balance{
		Model: model,
	}, nil
}

func (s *Service) GetTotalBalance(model *models.Address) (*Balance, error) {
	stakes, err := s.stakeRepository.GetAllByAddress(model.Address)
	if err != nil {
		return nil, err
	}

	// compute available balance from address balances
	totalBalanceSum := s.balanceService.GetTotalBalance(model)
	totalBalanceSumUSD := s.balanceService.GetTotalBalanceInUSD(totalBalanceSum)
	stakeBalanceSum := s.balanceService.GetStakeBalance(stakes)
	stakeBalanceSumUSD := s.balanceService.GetTotalBalanceInUSD(stakeBalanceSum)

	return &Balance{
		Model:              model,
		TotalBalanceSum:    totalBalanceSum,
		TotalBalanceSumUSD: totalBalanceSumUSD,
		StakeBalanceSum:    stakeBalanceSum,
		StakeBalanceSumUSD: stakeBalanceSumUSD,
	}, nil
}
