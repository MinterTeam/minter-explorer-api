package aggregated_reward

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/resource"
	"github.com/MinterTeam/minter-explorer-api/v2/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"time"
)

type Resource struct {
	TimeID    string             `json:"time_id"`
	Timestamp string             `json:"timestamp"`
	Role      string             `json:"role"`
	Amount    string             `json:"amount"`
	Address   string             `json:"address"`
	Validator resource.Interface `json:"validator"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	reward := model.(models.AggregatedReward)

	return Resource{
		TimeID:    reward.TimeID.Format(time.RFC3339),
		Timestamp: reward.TimeID.Format(time.RFC3339),
		Role:      reward.Role,
		Amount:    helpers.PipStr2Bip(reward.Amount),
		Address:   reward.Address.GetAddress(),
		Validator: new(validator.Resource).Transform(*reward.Validator),
	}
}
