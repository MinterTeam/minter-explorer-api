package slash

import (
	"github.com/MinterTeam/minter-explorer-api/events"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
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
		Column("Coin.symbol", "Address.address", "Block.created_at", "Validator").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return slashes
}

func (repository Repository) GetPaginatedByValidator(validator *models.Validator, pagination *tools.Pagination) ([]models.Slash, error) {
	var slashes []models.Slash
	var err error

	pagination.Total, err = repository.db.Model(&slashes).
		Column("Coin.symbol", "Address.address", "Block.created_at").
		Where("validator_id = ?", validator.ID).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	if err != nil {
		return nil, err
	}

	return slashes, nil
}
