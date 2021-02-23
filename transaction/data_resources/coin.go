package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"regexp"
	"strings"
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
	coin, err := coins.GlobalRepository.FindByID(uint(data.GetId()))
	helpers.CheckErr(err)

	return PoolCoin{
		ID:     data.GetId(),
		Symbol: coin.GetSymbol(),
		Amount: helpers.PipStr2Bip(pipAmount),
	}
}

func (PoolCoin) TransformCollection(data []*api_pb.Coin, model models.Transaction) []PoolCoin {
	pools := strings.Split(model.Tags["tx.pools"], ",")

	coins := make([]PoolCoin, len(data))
	for i, coin := range data {
		re := regexp.MustCompile(`-([0-9]*)`)
		amounts := re.FindAllString(pools[i/2], 2)

		amount := amounts[1][1:]
		if i == 0 {
			amount = amounts[0][1:]
		}

		coins[i] = new(PoolCoin).Transform(coin, amount)
	}

	return coins
}
