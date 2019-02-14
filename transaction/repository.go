package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type Repository struct {
	db *pg.DB
}

type SelectFilter struct {
	AddressId  *uint64
	StartBlock *string
	EndBlock   *string
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (f *SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.AddressId != nil {
		q = q.Join("LEFT OUTER JOIN transaction_outputs ON transaction_outputs.transaction_id = transaction.id").
			WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
				q = q.Where("transaction.from_address_id = ?", f.AddressId).
					WhereOr("transaction_outputs.to_address_id = ?", f.AddressId)

				return q, nil
			})
	}

	if f.StartBlock != nil {
		q = q.Where("transaction.block_id >= ?", f.StartBlock)
	}

	if f.EndBlock != nil {
		q = q.Where("transaction.block_id <= ?", f.EndBlock)
	}

	return q, nil
}

// Get paginated list of transactions by select filter
func (repository Repository) GetPaginatedTxByFilter(filter SelectFilter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Column("transaction.*", "FromAddress").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return transactions
}
