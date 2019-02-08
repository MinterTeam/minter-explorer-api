package coins

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type CoinResource struct {
	Crr                   uint64    `json:"crr"             example:"10"`
	Volume                string    `json:"volume"          example:"46573.556"`
	ReserveBalance        string    `json:"reserve_balance" example:"134.23456"`
	Name                  string    `json:"name"            example:"My test coin"`
	Symbol                string    `json:"symbol"          example:"TESTCOIN"`
}

func TransformCoin(model models.Coin) CoinResource  {
	return CoinResource{
		Crr:            model.Crr,
		Volume:         helpers.PipStr2Bip(model.Volume),
		ReserveBalance: helpers.PipStr2Bip(model.ReserveBalance),
		Name:           model.Name,
		Symbol:         model.Symbol,
	}
}