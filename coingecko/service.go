package coingecko

import "strings"

type CoinGeckoService struct {
	api *Api
}

func NewService(host string) *CoinGeckoService {
	return &CoinGeckoService{
		api: NewApi(host),
	}
}

func (s *CoinGeckoService) GetPrice(from string, to string) (float64, float64, error) {
	from, to = strings.ToLower(from), strings.ToLower(to)
	resp, err := s.api.GetSimplePrice(from, to)
	if err != nil {
		return 0, 0, err
	}

	return resp[from][to], resp[from][to+"_24h_change"], nil
}
