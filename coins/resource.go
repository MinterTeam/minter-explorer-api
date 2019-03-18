package coins

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	Crr            uint64 `json:"crr"`
	//Volume         string `json:"volume"` TODO: return when fields will be correct
	//ReserveBalance string `json:"reserveBalance"` TODO: return when fields will be correct
	Name           string `json:"name"`
	Symbol         string `json:"symbol"`
}

func (Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	coin := model.(models.Coin)
	return Resource{
		Crr:            coin.Crr,
		//Volume:         helpers.PipStr2Bip(coin.Volume),
		//ReserveBalance: helpers.PipStr2Bip(coin.ReserveBalance),
		Name:           coin.Name,
		Symbol:         coin.Symbol,
	}
}
