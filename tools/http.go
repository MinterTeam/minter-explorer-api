package tools

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	host   string
	client *http.Client
}

func NewHttpClient(host string) *HttpClient {
	return &HttpClient{
		host:   host,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *HttpClient) Get(url string, response interface{}) error {
	req, err := http.NewRequest("GET", c.host+"/"+url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return err
	}

	return nil
}
