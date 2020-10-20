package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
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
