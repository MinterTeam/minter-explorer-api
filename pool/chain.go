package pool

import "encoding/json"

type SwapChain struct {
	PoolId   uint64 `json:"pool_id"`
	CoinIn   uint64 `json:"coin_in"`
	ValueIn  string `json:"value_in"`
	CoinOut  uint64 `json:"coin_out"`
	ValueOut string `json:"value_out"`
}

func GetPoolChainFromStr(chainStr string) (chain []SwapChain) {
	json.Unmarshal([]byte(chainStr), &chain)
	return chain
}
