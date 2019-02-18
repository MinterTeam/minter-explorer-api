package slash

import (
	"github.com/MinterTeam/minter-explorer-api/events"
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

func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Slash {
	var slashes []models.Slash
	var err error

	pagination.Total, err = repository.db.Model(&slashes).
		Column("Coin.symbol", "Address.address", "Validator.public_key").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return slashes
}