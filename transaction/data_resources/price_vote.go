package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type PriceVote struct {
	Price string `json:"price"`
}

func (PriceVote) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.PriceVoteData)

	return PriceVote{
		Price: data.Price,
	}
}
