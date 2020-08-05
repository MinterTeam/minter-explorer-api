package services

import (
	"github.com/MinterTeam/minter-explorer-api/stake"
	"github.com/MinterTeam/minter-explorer-api/tools/cache"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/sirupsen/logrus"
)

type ValidatorService struct {
	stakeRepository *stake.Repository
	repository      *validator.Repository
	cache           *cache.ExplorerCache
}

func NewValidatorService(
	repository *validator.Repository,
	stakeRepository *stake.Repository,
	cache *cache.ExplorerCache,
) *ValidatorService {
	return &ValidatorService{
		stakeRepository: stakeRepository,
		repository:      repository,
		cache:           cache,
	}
}

const (
	minStakeCacheTime = 120
)

type ValidatorsMinStake map[uint64]string

func (s *ValidatorService) GetMinStakesByValidator(validatorId uint64) string {
	var stakes ValidatorsMinStake

	if s.cache.GetLastBlock().ID%120 == 0 {
		stakes = s.getMinStakes()
		s.cache.Store("min_stakes", stakes, minStakeCacheTime)
	} else {
		stakes = s.cache.Get("min_stakes", func() interface{} {
			return s.getMinStakes()
		}, minStakeCacheTime).(ValidatorsMinStake)
	}

	if bipStake, ok := stakes[validatorId]; ok {
		return bipStake
	}

	return "0"
}

func (s *ValidatorService) getMinStakes() ValidatorsMinStake {
	stakes, err := s.stakeRepository.GetMinStakes()
	if err != nil {
		logrus.WithField("err", err).Error("min stakes error")
	}

	minStakes := make(ValidatorsMinStake, len(stakes))
	for _, v := range stakes {
		minStakes[v.ValidatorID] = v.BipValue
	}

	return minStakes
}
