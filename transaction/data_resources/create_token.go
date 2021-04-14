package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
)

type CreateToken struct {
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	InitialAmount string `json:"initial_amount"`
	MaxSupply     string `json:"max_supply"`
	Mintable      bool   `json:"mintable"`
	Burnable      bool   `json:"burnable"`
}

func (CreateToken) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.CreateTokenData)

	return CreateToken{
		Name:          data.Name,
		Symbol:        data.Symbol,
		Mintable:      data.Mintable,
		Burnable:      data.Burnable,
		InitialAmount: helpers.PipStr2Bip(data.InitialAmount),
		MaxSupply:     helpers.PipStr2Bip(data.MaxSupply),
	}
}
