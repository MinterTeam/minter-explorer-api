package coins

import (
	"encoding/json"
	"github.com/MinterTeam/minter-explorer-api/v2/hub"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Service struct {
	hubClient     *hub.Client
	repository    *Repository
	verifiedCoins []models.Coin
	blockListCoinIds []uint64
}

func NewService(repository *Repository) *Service {
	s := &Service{
		hubClient:  hub.NewClient(),
		repository: repository,
	}

	go s.RunVerifiedCoinsUpdater()
	go s.RunBlocklistCoinsUpdater()

	return s
}

func (s *Service) GetVerifiedCoins() []models.Coin {
	if len(s.verifiedCoins) == 0 {
		if verifiedCoins, err := s.getVerifiedCoins(); err == nil {
			return verifiedCoins
		}
	}

	return s.verifiedCoins
}

func (s *Service) GetBlocklistCoins() []uint64 {
	if len(s.blockListCoinIds) == 0 {
		if coinIds, err := s.getBlocklistCoinIds(); err == nil {
			return coinIds
		}
	}

	return s.blockListCoinIds
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

func (s *Service) RunBlocklistCoinsUpdater() {
	for {
		coins, err := s.getBlocklistCoinIds()
		if err != nil {
			log.Errorf("failed to get blocklist coins: %s", err)
			continue
		}

		s.blockListCoinIds = coins
		s.repository.SetBlocklistCoinIds(s.blockListCoinIds)
		time.Sleep(time.Minute * 15)
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

func (s *Service) getBlocklistCoinIds() ([]uint64, error) {
	resp, err := resty.New().R().Get("https://raw.githubusercontent.com/MinterTeam/coin-block-list/master/dist/blocklist.json")
	if err != nil || resp.IsError() {
		return nil, err
	}

	var coinsSymbol []string
	err = json.Unmarshal(resp.Body(), &coinsSymbol)
	if err != nil {
		return nil, err
	}

	coins, err := s.repository.GetBySymbols(coinsSymbol)
	if err != nil {
		return nil, err
	}

	coinIds := make([]uint64, len(coins))
	for i, c := range coins {
		coinIds[i] = uint64(c.ID)
	}

	return coinIds, nil
}