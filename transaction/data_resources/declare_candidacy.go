package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type DeclareCandidacy struct {
	Address    string `json:"address"`
	PubKey     string `json:"pub_key"`
	Commission string `json:"commission"`
	Coin       Coin   `json:"coin"`
	Stake      string `json:"stake"`
}

func (DeclareCandidacy) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.DeclareCandidacyData)

	return DeclareCandidacy{
		Address:    data.Address,
		PubKey:     data.PubKey,
		Commission: data.Commission,
		Coin:       new(Coin).Transform(data.Coin),
		Stake:      helpers.PipStr2Bip(data.Stake),
	}
}
