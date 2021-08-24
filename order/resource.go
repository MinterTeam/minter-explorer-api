package order

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Address          string             `json:"address"`
	PoolId           uint64             `json:"pool_id"`
	CoinToSell       resource.Interface `json:"coin_to_sell"`
	CoinToBuy        resource.Interface `json:"coin_to_buy"`
	CoinToBuyVolume  string             `json:"coin_to_buy_volume"`
	CoinToSellVolume string             `json:"coin_to_sell_volume"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	order := model.(models.Order)
	return Resource{
		Address:          order.Address.GetAddress(),
		PoolId:           order.LiquidityPoolId,
		CoinToSell:       nil,
		CoinToBuy:        nil,
		CoinToBuyVolume:  helpers.PipStr2Bip(order.CoinBuyVolume),
		CoinToSellVolume: helpers.PipStr2Bip(order.CoinSellVolume),
	}
}
