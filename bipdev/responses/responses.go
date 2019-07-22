package responses

import "time"

type CurrentPriceResponse struct {
	Data struct {
		Price      float64   `json:"price"`
		NextUpdate time.Time `json:"next_update"`
		Delta      float64   `json:"delta"`
	} `json:"data"`
}
