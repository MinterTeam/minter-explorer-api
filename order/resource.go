package order

import (
	"github.com/MinterTeam/minter-explorer-api/v2/coins"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/node-grpc-gateway/api_pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type Resource struct {
	Id                      uint64             `json:"id"`
	Address                 string             `json:"address"`
	PoolId                  uint64             `json:"pool_id"`
	InitialCoinToBuyVolume  string             `json:"initial_coin_to_buy_volume"`
	CoinToBuyVolume         string             `json:"coin_to_buy_volume"`
	InitialCoinToSellVolume string             `json:"initial_coin_to_sell_volume"`
	CoinToSellVolume        string             `json:"coin_to_sell_volume"`
	CoinToSell              resource.Interface `json:"coin_to_sell"`
	CoinToBuy               resource.Interface `json:"coin_to_buy"`
	Height                  uint64             `json:"height"`
	Status                  Status             `json:"status"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	order := model.(OrderTransaction)

	txData := reflect.New(reflect.TypeOf(new(api_pb.AddLimitOrderData)).Elem()).Interface().(proto.Message)
	protojson.Unmarshal(order.Transaction.Data, txData)
	limitOrder := txData.(*api_pb.AddLimitOrderData)

	return Resource{
		Id:                      order.Id,
		Address:                 order.Address.GetAddress(),
		PoolId:                  order.LiquidityPoolId,
		CoinToSell:              new(coins.IdResource).Transform(*order.CoinSell),
		CoinToBuy:               new(coins.IdResource).Transform(*order.CoinBuy),
		CoinToBuyVolume:         helpers.PipStr2Bip(order.CoinBuyVolume),
		CoinToSellVolume:        helpers.PipStr2Bip(order.CoinSellVolume),
		InitialCoinToBuyVolume:  helpers.PipStr2Bip(limitOrder.ValueToBuy),
		InitialCoinToSellVolume: helpers.PipStr2Bip(limitOrder.ValueToSell),
		Height:                  order.CreatedAtBlock,
		Status:                  MakeOrderStatus(order.Status),
	}
}
