package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"strconv"
)

type Coin struct {
	ID     uint32 `json:"id"`
	Symbol string `json:"symbol"`
}

func (Coin) Transform(data *api_pb.Coin) Coin {
	id, _ := strconv.ParseUint(data.GetId(), 10, 32)

	coin, err := coins.GlobalRepository.FindByID(uint(id))
	helpers.CheckErr(err)
	
	return Coin{
		ID:     uint32(id),
		Symbol: coin.GetSymbol(),
	}
}
