package coins

import (
	"github.com/MinterTeam/minter-explorer-api/v2/hub"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Service struct {
	hubClient     *hub.Client
	repository    *Repository
	verifiedCoins []models.Coin
}

func NewService(repository *Repository) *Service {
	return &Service{
		hubClient:  hub.NewClient(),
		repository: repository,
	}
}

func (s *Service) GetVerifiedCoins() []models.Coin {
	if len(s.verifiedCoins) == 0 {
		if verifiedCoins, err := s.getVerifiedCoins(); err == nil {
			return verifiedCoins
		}
	}

	return s.verifiedCoins
}

func (s *Service) RunVerifiedCoinsUpdater() {
	for {
		coins, err := s.getVerifiedCoins()
		if err != nil {
			log.Errorf("failed to get verified coins: %s", err)
			continue
		}

		s.verifiedCoins = coins
		time.Sleep(time.Minute * 5)
	}
}

func (s *Service) getVerifiedCoins() ([]models.Coin, error) {
	resp, err := s.hubClient.GetOracleCoins()
	if err != nil {
		return nil, err
	}

	coinIds := make([]models.Coin, len(resp.Result))
	for i, coin := range resp.Result {
		coinId, _ := strconv.ParseUint(coin.MinterId, 10, 64)
		coinIds[i], _ = s.repository.FindByID(uint(coinId))
	}

	return coinIds, nil
}
