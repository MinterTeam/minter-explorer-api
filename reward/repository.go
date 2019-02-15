package reward

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repository Repository) GetPaginatedByAddress(address string, pagination *tools.Pagination) []models.Reward {
	var rewards []models.Reward
	var err error

	pagination.Total, err = repository.db.Model(&rewards).
		Column("Address", "Validator").
		Where("address.address = ?", address).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return rewards
}