package waitlist

import (
	"github.com/MinterTeam/minter-explorer-api/coins"
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-api/validator"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
)

type Resource struct {
	Value     string             `json:"value"`
	Coin      resource.Interface `json:"coin"`
	Validator resource.Interface `json:"validator"`
}

func (Resource) Transform(model resource.ItemInterface, params ...resource.ParamInterface) resource.Interface {
	wl := model.(models.StakeKick)

	return Resource{
		Value:     wl.Value,
		Coin:      new(coins.IdResource).Transform(*wl.Coin),
		Validator: new(validator.Resource).Transform(*wl.Validator),
	}
}
