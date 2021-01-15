package check

import (
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v9"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetByRawCheck(raw string) (models.Check, error) {
	var check models.Check

	err := r.db.Model(&check).
		Relation("FromAddress").
		Relation("ToAddress").
		Where("data = ?", raw).
		Select()

	return check, err
}

func (r *Repository) GetListByFilter(filter SelectFilter, pagination *tools.Pagination) (checks []models.Check, err error) {
	pagination.Total, err = r.db.Model(&checks).
		Relation("FromAddress").
		Relation("ToAddress").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		SelectAndCount()

	return checks, err
}
