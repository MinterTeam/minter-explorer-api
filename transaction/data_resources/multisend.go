package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"strconv"
)

type Multisend struct {
	List []Send `json:"list"`
}

func (Multisend) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	data := txData.(*api_pb.MultiSendData)

	list := make([]Send, len(data.List))
	for key, item := range data.List {
		list[key] = Send{}.Transform(item).(Send)
	}

	return Multisend{list}
}

func (Multisend) TransformByTxOutput(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.TransactionOutput)
	coin := &api_pb.Coin{
		Id:     strconv.FormatUint(uint64(data.Coin.ID), 10),
		Symbol: data.Coin.Symbol,
	}

	return Send{
		Coin:  new(Coin).Transform(coin),
		To:    data.ToAddress.GetAddress(),
		Value: helpers.PipStr2Bip(data.Value),
	}
}
