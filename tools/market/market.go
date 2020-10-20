package market

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
	"time"
)

type Service struct {
	api         ApiService
	PriceChange PriceChange
	basecoin    string
}

type ApiService interface {
	GetPrice(from string, to string) (float64, float64, error)
}

type PriceChange struct {
	Price  float64
	Change float64
}

func NewService(api ApiService, basecoin string) *Service {
	return &Service{
		api:         api,
		basecoin:    basecoin,
		PriceChange: PriceChange{Price: 0, Change: 0},
	}
}

func (s *Service) Run() {
	for {
		price, change, err := s.api.GetPrice(s.basecoin, USDTicker)
		if err == nil {
			s.PriceChange = PriceChange{
				Price:  price,
				Change: change,
			}
		}

		time.Sleep(config.MarketPriceUpdatePeriodInMin * time.Minute)
	}
}
