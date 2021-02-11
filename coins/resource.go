package coins

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	ID           uint    `json:"id"`
	Crr          uint    `json:"crr"`
	Volume       string  `json:"volume"`
	Reserve      string  `json:"reserve_balance"`
	MaxSupply    string  `json:"max_supply"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	OwnerAddress *string `json:"owner_address"`
	Burnable     bool    `json:"burnable"`
	Mintable     bool    `json:"mintable"`
	Type         string  `json:"type"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	coin := model.(models.Coin)

	ownerAddress := new(string)
	if coin.OwnerAddressId != 0 {
		*ownerAddress = coin.OwnerAddress.GetAddress()
	} else {
		ownerAddress = nil
	}

	return Resource{
		ID:           coin.ID,
		Crr:          coin.Crr,
		Volume:       helpers.PipStr2Bip(coin.Volume),
		Reserve:      helpers.PipStr2Bip(coin.Reserve),
		MaxSupply:    helpers.PipStr2Bip(coin.MaxSupply),
		Name:         coin.Name,
		Symbol:       coin.GetSymbol(),
		OwnerAddress: ownerAddress,
		Burnable:     coin.Burnable,
		Mintable:     coin.Mintable,
		Type:         helpers.GetCoinType(coin.Type),
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
		Symbol: coin.GetSymbol(),
	}
}
