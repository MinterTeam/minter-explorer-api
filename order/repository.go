package order

import (
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

func (r *Repository) GetListPaginated(pagination *tools.Pagination, filters ...tools.Filter) (orders []models.Order, err error) {
	q := r.db.Model(&orders).
		Relation("Address").
		Relation("CoinSell").
		Relation("CoinBuy").
		Order("coin_sell_volume / coin_buy_volume desc").
		Apply(pagination.Filter)

	for _, f := range filters {
		q = q.Apply(f.Filter)
	}

	pagination.Total, err = q.SelectAndCount()
	return
}
