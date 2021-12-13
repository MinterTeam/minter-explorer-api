package hub

import "github.com/go-resty/resty/v2"

type Client struct {
	http *resty.Client
}

func NewClient() *Client {
	return &Client{
		http: resty.New().SetHostURL("https://hub-api.minter.network"),
	}
}

type OracleCoinsResponse struct {
	List struct {
		TokenInfos []struct {
			Id               string `json:"id"`
			Denom            string `json:"denom"`
			ChainId          string `json:"chain_id"`
			ExternalTokenId  string `json:"external_token_id"`
			ExternalDecimals string `json:"external_decimals"`
			Commission       string `json:"commission"`
		} `json:"token_infos"`
	} `json:"list"`
}

func (c *Client) GetOracleCoins() (*OracleCoinsResponse, error) {
	resp := new(OracleCoinsResponse)
	_, err := c.http.R().SetResult(resp).Get("/mhub2/v1/token_infos")
	return resp, err
}
