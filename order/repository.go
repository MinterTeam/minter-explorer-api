package order

import (
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-go-sdk/v2/transaction"
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
		Column("coin_sell_volume", "coin_buy_volume", "created_at_block", "status", "liquidity_pool_id", "price").
		ColumnExpr(`"order_transaction".id AS "id"`).
		ColumnExpr(`transactions.data AS "transaction__data"`).
		Join(`JOIN transactions ON (transactions.tags->>'tx.order_id')::int = "order_transaction".id and transactions.type = ?`, transaction.TypeAddLimitOrder).
		Apply(pagination.Filter)

	for _, f := range filters {
		q = q.Apply(f.Filter)
	}

	q = q.Order("id DESC")

	pagination.Total, err = q.SelectAndCount()
	return
}

func (r *Repository) FindById(id uint64) (OrderTransaction, error) {
	var order OrderTransaction
	err := r.db.Model(&order).
		Relation("Address").
		Relation("CoinSell").
		Relation("CoinBuy").
		Column("coin_sell_volume", "coin_buy_volume", "created_at_block", "status", "liquidity_pool_id", "price").
		ColumnExpr(`"order_transaction".id AS "id"`).
		ColumnExpr(`transactions.data AS "transaction__data"`).
		Join(`JOIN transactions ON (transactions.tags->>'tx.order_id')::int = "order_transaction".id and transactions.type = ?`, transaction.TypeAddLimitOrder).
		Where(`"order_transaction".id = ?`, id).
		Select()
	return order, err
}
