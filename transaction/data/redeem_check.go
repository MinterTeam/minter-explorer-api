package data

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type RedeemCheckResource struct {
	RawCheck string `json:"raw_check"`
	Proof    string `json:"proof"`
}

func (RedeemCheckResource) Transform(txData resource.ItemInterface) resource.Interface {
	data := txData.(*models.RedeemCheckTxData)

	return RedeemCheckResource{
		RawCheck: data.RawCheck,
		Proof:    data.Proof,
	}
}
