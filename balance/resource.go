package balance

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/MinterTeam/minter-go-node/formula"
	"math/big"
)

type Resource struct {
	Coin      resource.Interface `json:"coin"`
	Amount    string             `json:"amount"`
	BipAmount string             `json:"bip_amount"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	balance := model.(models.Balance)

	return Resource{
		Coin:      new(coins.IdResource).Transform(*balance.Coin),
		Amount:    helpers.PipStr2Bip(balance.Value),
		BipAmount: helpers.PipStr2Bip(getCoinBalanceInBaseValue(balance).String()),
	}
}

func getCoinBalanceInBaseValue(balance models.Balance) *big.Int {
	if balance.Coin.ID == 0 {
		return helpers.StringToBigInt(balance.Value)
	}

	return formula.CalculateSaleReturn(
		helpers.StringToBigInt(balance.Coin.Volume),
		helpers.StringToBigInt(balance.Coin.Reserve),
		balance.Coin.Crr,
		helpers.StringToBigInt(balance.Value),
	)
}
