package chart

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/reward"
	"time"
)

type RewardResource struct {
	Time   string `json:"time"`
	Amount string `json:"amount"`
}

func (RewardResource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	data := model.(reward.ChartData)

	return RewardResource{
		Time:   data.Time.Format(time.RFC3339),
		Amount: helpers.PipStr2Bip(data.Amount),
	}
}
