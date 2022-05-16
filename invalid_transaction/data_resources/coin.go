package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
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
}

func (PoolCoin) Transform(data *api_pb.Coin) PoolCoin {
	return PoolCoin{
		ID:     data.GetId(),
		Symbol: data.GetSymbol(),
	}
}


func (PoolCoin) TransformCollection(data []*api_pb.Coin, model models.InvalidTransaction) []PoolCoin {
	poolCoins := make([]PoolCoin, len(data))

	for i, coin := range data {
		// todo: refactor
		coinModel, err := coins.GlobalRepository.FindByID(uint(coin.Id))
		helpers.CheckErr(err)
		coin.Symbol = coinModel.GetSymbol()
		poolCoins[i] = new(PoolCoin).Transform(coin)
	}

	return poolCoins
}
