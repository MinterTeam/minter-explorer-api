package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type RedeemCheck struct {
	RawCheck string `json:"raw_check"`
	Proof    string `json:"proof"`
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...interface{}) resource.Interface {
	data := txData.(*models.RedeemCheckTxData)

	return RedeemCheck{
		RawCheck: data.RawCheck,
		Proof:    data.Proof,
	}
}
