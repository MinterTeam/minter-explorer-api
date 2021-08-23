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
	Height string `json:"height"`
	Result []struct {
		Denom            string `json:"denom"`
		EthAddr          string `json:"eth_addr"`
		MinterId         string `json:"minter_id"`
		EthDecimals      string `json:"eth_decimals"`
		CustomCommission string `json:"custom_commission"`
	} `json:"result"`
}

func (c *Client) GetOracleCoins() (*OracleCoinsResponse, error) {
	resp := new(OracleCoinsResponse)
	_, err := c.http.R().SetResult(resp).Get("/oracle/coins")
	return resp, err
}
