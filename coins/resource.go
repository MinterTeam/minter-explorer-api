package coins

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Crr          uint    `json:"crr"`
	Volume       string  `json:"volume"`
	Reserve      string  `json:"reserve_balance"`
	MaxSupply    string  `json:"max_supply"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	OwnerAddress *string `json:"owner_address"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	ownerAddress := new(string)
	if coin.OwnerAddressId != 0 {
		ownerAddress = &coin.OwnerAddress.Address
	}

	return Resource{
		Crr:          coin.Crr,
		Volume:       helpers.PipStr2Bip(coin.Volume),
		Reserve:      helpers.PipStr2Bip(coin.Reserve),
		MaxSupply:    helpers.PipStr2Bip(coin.MaxSupply),
		Name:         coin.Name,
		Symbol:       coin.Symbol,
		OwnerAddress: ownerAddress,
	}
}

type IdResource struct {
	ID     uint32 `json:"id"`
	Symbol string `json:"symbol"`
}

func (IdResource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	return IdResource{
		ID:     uint32(coin.ID),
		Symbol: coin.Symbol,
	}
}
