package meta

import (
	"github.com/MinterTeam/minter-explorer-api/resource"
	"github.com/MinterTeam/minter-explorer-tools/models"
)

type Resource struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IconUrl     *string `json:"icon_url"`
	SiteUrl     *string `json:"site_url"`
}

func (r Resource) Transform(model resource.ItemInterface, params ...interface{}) resource.Interface {
	validator := model.(models.Validator)

	return Resource{
		Name:        validator.Name,
		Description: validator.Description,
		IconUrl:     validator.IconUrl,
		SiteUrl:     validator.SiteUrl,
	}
}
