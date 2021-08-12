package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/pool"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type Coin struct {
	ID     uint64 `json:"id"`
	Symbol string `json:"symbol"`
}

func (Coin) Transform(data *api_pb.Coin) Coin {
	coin, err := coins.GlobalRepository.FindByID(uint(data.GetId()))
	helpers.CheckErr(err)

	return Coin{
		ID:     data.GetId(),
		Symbol: coin.GetSymbol(),
	}
}

type PoolCoin struct {
	ID     uint64 `json:"id"`
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

func (PoolCoin) Transform(data *api_pb.Coin, pipAmount string) PoolCoin {
	return PoolCoin{
		ID:     data.GetId(),
		Symbol: data.GetSymbol(),
		Amount: helpers.PipStr2Bip(pipAmount),
	}
}

func (PoolCoin) TransformCollection(data []*api_pb.Coin, model models.Transaction) []PoolCoin {
	pools := pool.GetPoolChainFromStr(model.Tags["tx.pools"])
	poolCoins := make([]PoolCoin, len(data))

	for i, coin := range data {
		pool := pools[0]
		if len(pools) > 1 {
			if i == len(pools) {
				pool = pools[i-1]
			} else {
				pool = pools[i]
			}
		}

		amount := pool.ValueIn
		if i == len(pools) {
			amount = pool.ValueOut
		}

		// todo: refactor
		coinModel, err := coins.GlobalRepository.FindByID(uint(coin.Id))
		helpers.CheckErr(err)
		coin.Symbol = coinModel.GetSymbol()

		poolCoins[i] = new(PoolCoin).Transform(coin, amount)
	}

	return poolCoins
}
