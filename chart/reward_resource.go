package chart

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/reward"
)

type RewardResource struct {
	Time   string `json:"time"`
	Amount string `json:"amount"`
}

func (RewardResource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	data := model.(reward.ChartData)

	return RewardResource{
		Time:   data.Time,
		Amount: helpers.PipStr2Bip(data.Amount),
	}
}
