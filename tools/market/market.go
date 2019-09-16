package market

import (
	"github.com/MinterTeam/minter-explorer-api/bipdev"
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"time"
)

type Service struct {
	PriceChange PriceChange
	api         *bipdev.Api
	baseCoin    string
}

type PriceChange struct {
	Price  float64
	Change float64
}

func NewService(bipdevApi *bipdev.Api, basecoin string) *Service {
	return &Service{
		api:         bipdevApi,
		baseCoin:    basecoin,
		PriceChange: PriceChange{Price: 0, Change: 0},
	}
}

func (s *Service) Run() {
	for {
		response, err := s.api.GetCurrentPrice()
		if err == nil {
			s.PriceChange = PriceChange{
				Price:  response.Data.Price / 10000,
				Change: response.Data.Delta,
			}
		}

		time.Sleep(time.Duration(config.MarketPriceUpdatePeriodInMin * time.Minute))
	}
}
