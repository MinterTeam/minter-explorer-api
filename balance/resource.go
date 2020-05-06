package balance

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/v4/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Resource struct {
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	BipAmount string `json:"bip_amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)

	return Resource{
		Coin:      balance.Coin.Symbol,
		Amount:    helpers.PipStr2Bip(balance.Value),
		BipAmount: helpers.PipStr2Bip(getCoinBalanceInBaseValue(balance).String()),
	}
}

func getCoinBalanceInBaseValue(balance models.Balance) *big.Int {
	if balance.Coin.ReserveBalance == "" {
		return helpers.StringToBigInt(balance.Value)
	}

	return formula.CalculateSaleReturn(
		helpers.StringToBigInt(balance.Coin.Volume),
		helpers.StringToBigInt(balance.Coin.ReserveBalance),
		uint(balance.Coin.Crr),
		helpers.StringToBigInt(balance.Value),
	)
}
