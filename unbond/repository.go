package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/events"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetListByAddress(filter *events.SelectFilter, pagination *tools.Pagination) ([]models.Unbond, error) {
	var unbonds []models.Unbond
	var err error

	pagination.Total, err = r.db.Model(&unbonds).
		Join("JOIN addresses as address ON address.id = unbond.address_id").
		Apply(filter.Filter).
		SelectAndCount()

	return unbonds, err
}
