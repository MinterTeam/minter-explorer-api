package coingecko

import (
	"fmt"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
)

type Api struct {
	client *tools.HttpClient
}

func NewApi(host string) *Api {
	return &Api{
		client: tools.NewHttpClient(host),
	}
}

type tickerData map[string]float64
type SimplePriceResponse map[string]tickerData

func (api *Api) GetSimplePrice(from string, to string) (SimplePriceResponse, error) {
	response := make(SimplePriceResponse)
	url := fmt.Sprintf("/api/v3/simple/price?include_24hr_change=true&ids=%s&vs_currencies=%s", from, to)
	if err := api.client.Get(url, &response); err != nil {
		return nil, err
	}

	return response, nil
}
