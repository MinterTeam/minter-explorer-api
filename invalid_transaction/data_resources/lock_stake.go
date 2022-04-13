package data_resources

import (
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
)

type LockStakeData struct{}

func (LockStakeData) Transform(txData resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	//data := txData.(*api_pb.LockStakeData)
	return LockStakeData{}
}
