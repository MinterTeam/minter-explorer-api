package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
)

type RedeemCheck struct {
}

func (RedeemCheck) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	return RedeemCheck{}
}
