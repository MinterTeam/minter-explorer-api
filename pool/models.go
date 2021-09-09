package pool

import (
	"github.com/MinterTeam/explorer-sdk/swap"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type SwapChain struct {
	PoolId   int    `json:"pool_id"`
	CoinIn   int    `json:"coin_in"`
	ValueIn  string `json:"value_in"`
	CoinOut  int    `json:"coin_out"`
	ValueOut string `json:"value_out"`
	Details  struct {
		AmountIn            string `json:"amount_in"`
		AmountOut           string `json:"amount_out"`
		CommissionAmountIn  string `json:"commission_amount_in"`
		CommissionAmountOut string `json:"commission_amount_out"`
		Orders              []struct {
			Buy    string `json:"buy"`
			Sell   string `json:"sell"`
			Seller string `json:"seller"`
			Id     int    `json:"id"`
		} `json:"orders"`
	} `json:"details"`
	Sellers []struct {
		Seller string `json:"seller"`
		Value  string `json:"value"`
	} `json:"sellers"`
}

type tradeVolume struct {
	PoolId           uint64
	Date             time.Time
	FirstCoinVolume  string
	SecondCoinVolume string
}

type TradeVolume struct {
	Date             time.Time
	FirstCoinVolume  string
	SecondCoinVolume string
	BipVolume        *big.Float
}

type TradeVolumes struct {
	Day   TradeVolume
	Month TradeVolume
}

type TradeSearch struct {
	Log        *log.Entry
	FromCoinId uint64
	ToCoinId   uint64
	TradeType  swap.TradeType
	Amount     *big.Int
	Trade      chan *swap.Trade
}
