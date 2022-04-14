package services

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"github.com/MinterTeam/minter-explorer-api/v2/stake"
	"github.com/MinterTeam/minter-explorer-api/v2/tools/cache"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	log "github.com/sirupsen/logrus"
)

type ValidatorService struct {
	stakeRepo *stake.Repository
	cache     *cache.ExplorerCache
}

type StakeRepository interface {
	GetMinStakes() ([]models.Stake, error)
}

func NewValidatorService(stakeRepo *stake.Repository, cache *cache.ExplorerCache) *ValidatorService {
	return &ValidatorService{
		stakeRepo: stakeRepo,
		cache:     cache,
	}
}

const (
	minStakeCacheTime = 720
)

type ValidatorsMinStake map[uint]string

func (s *ValidatorService) GetMinStakesByValidator(validator *models.Validator) string {
	if len(validator.Stakes) < config.MaxDelegatorCount {
		return "0"
	}

	var stakes ValidatorsMinStake

	if s.cache.GetLastBlock().ID%120 == 0 {
		stakes = s.getMinStakes()
		s.cache.Store("min_stakes", stakes, minStakeCacheTime)
	} else {
		stakes = s.cache.Get("min_stakes", func() interface{} {
			return s.getMinStakes()
		}, minStakeCacheTime).(ValidatorsMinStake)
	}

	if bipStake, ok := stakes[validator.ID]; ok {
		return bipStake
	}

	return "0"
}

func (s *ValidatorService) getMinStakes() ValidatorsMinStake {
	stakes, err := s.stakeRepo.GetMinStakes()
	if err != nil {
		log.Errorf("failed to get min stakes: %s", err)
		return nil
	}

	minStakes := make(ValidatorsMinStake, len(stakes))
	for _, v := range stakes {
		minStakes[v.ValidatorID] = v.BipValue
	}

	return minStakes
}
