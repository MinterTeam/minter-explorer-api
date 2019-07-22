package market

import (
	"github.com/MinterTeam/minter-explorer-api/bipdev"
)

type Service struct {
	api      *bipdev.Api
	baseCoin string
}

func NewService(bipdevApi *bipdev.Api, basecoin string) *Service {
	return &Service{
		api:      bipdevApi,
		baseCoin: basecoin,
	}
}

type PriceChange struct {
	Price  float64
	Change float64
}

func (s *Service) GetCurrentFiatPriceChange(coin string, currency string) (*PriceChange, error) {
	if coin == s.baseCoin && currency == USDTicker {
		response, err := s.api.GetCurrentPrice()
		if err != nil {
			return nil, err
		}

		return &PriceChange{
			Price:  response.Data.Price / 10000,
			Change: response.Data.Delta,
		}, nil
	}

	return nil, nil
}
