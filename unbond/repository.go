package unbond

import (
	"github.com/MinterTeam/minter-explorer-api/v2/events"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
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
		Relation("Coin").
		Relation("Validator").
		ColumnExpr("unbond.block_id, unbond.value, address.address as address__address").
		Join("JOIN addresses as address ON address.id = unbond.address_id").
		Apply(filter.Filter).
		SelectAndCount()

	return unbonds, err
}
