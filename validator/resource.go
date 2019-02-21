package validator

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/stake"
	"github.com/MinterTeam/minter-explorer-extender/models"
)

type Resource struct {
	Status         uint8                `json:"status"`
	Stake          string               `json:"stake"`
	Part           string               `json:"part"`
	DelegatorCount int                  `json:"delegator_count"`
	DelegatorList  []resource.Interface `json:"delegator_list"`
	TotalStake     string               `json:"-"`
}

func (r Resource) Transform(model resource.ItemInterface) resource.Interface {
	validator := model.(models.Validator)
	delegators := resource.TransformCollection(validator.Stakes, stake.Resource{})

	return Resource{
		Status:         *validator.Status,
		Stake:          helpers.PipStr2Bip(*validator.TotalStake),
		Part:           helpers.CalculatePercent(*validator.TotalStake, r.TotalStake),
		DelegatorCount: len(delegators),
		DelegatorList:  delegators,
	}
}
