package pool

import (
	"github.com/MinterTeam/explorer-sdk/swap"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type SwapChain struct {
	PoolId   uint64 `json:"pool_id"`
	CoinIn   uint64 `json:"coin_in"`
	ValueIn  string `json:"value_in"`
	CoinOut  uint64 `json:"coin_out"`
	ValueOut string `json:"value_out"`
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
