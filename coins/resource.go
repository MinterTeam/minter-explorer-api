package coins

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Resource struct {
	Crr            uint64 `json:"crr"             example:"10"`
	Volume         string `json:"volume"          example:"46573.556"`
	ReserveBalance string `json:"reserve_balance" example:"134.23456"`
	Name           string `json:"name"            example:"My test coin"`
	Symbol         string `json:"symbol"          example:"TESTCOIN"`
}

func (Resource) Transform(model resource.ItemInterface) resource.Interface {
	coin := model.(models.Coin)
	return Resource{
		Crr:            coin.Crr,
		Volume:         helpers.PipStr2Bip(coin.Volume.String()),
		ReserveBalance: helpers.PipStr2Bip(coin.ReserveBalance.String()),
		Name:           coin.Name,
		Symbol:         coin.Symbol,
	}
}
