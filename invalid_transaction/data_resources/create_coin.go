package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type CreateCoin struct {
	Name                 string `json:"name"`
	Symbol               string `json:"symbol"`
	InitialAmount        string `json:"initial_amount"`
	InitialReserve       string `json:"initial_reserve"`
	ConstantReserveRatio uint64 `json:"constant_reserve_ratio"`
	MaxSupply            string `json:"max_supply"`
}

func (CreateCoin) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.CreateCoinData)

	return CreateCoin{
		Name:                 data.Name,
		Symbol:               data.Symbol,
		InitialAmount:        helpers.PipStr2Bip(data.InitialAmount),
		InitialReserve:       helpers.PipStr2Bip(data.InitialReserve),
		ConstantReserveRatio: data.ConstantReserveRatio,
		MaxSupply:            helpers.PipStr2Bip(data.MaxSupply),
	}
}
