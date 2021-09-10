package order

import (
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/go-pg/pg/v10"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetListPaginated(pagination *tools.Pagination, filters ...tools.Filter) (orders []OrderTransaction, err error) {
	q := r.db.Model(&orders).
		Relation("Address").
		Relation("CoinSell").
		Relation("CoinBuy").
		Column("coin_sell_volume", "coin_buy_volume", "created_at_block", "status", "liquidity_pool_id").
		ColumnExpr(`"order_transaction".id AS "id"`).
		ColumnExpr(`transactions.data AS "transaction__data"`).
		Join(`JOIN transactions ON (transactions.tags->>'tx.order_id')::int = "order_transaction".id`).
		OrderExpr("coin_sell_volume / coin_buy_volume desc").
		Apply(pagination.Filter)

	for _, f := range filters {
		q = q.Apply(f.Filter)
	}

	pagination.Total, err = q.SelectAndCount()
	return
}
