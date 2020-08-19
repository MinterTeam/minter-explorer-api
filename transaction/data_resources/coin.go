package data_resources

import (
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"strconv"
)

type Coin struct {
	ID     uint32 `json:"id"`
	Symbol string `json:"symbol"`
}

func (Coin) Transform(data *api_pb.Coin) Coin {
	id, _ := strconv.ParseUint(data.GetId(), 10, 32)

	return Coin{
		ID:     uint32(id),
		Symbol: data.GetSymbol(),
	}
}
