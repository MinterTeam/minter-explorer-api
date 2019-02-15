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
	Addresses       []interface{}
	BlockId         *uint64
	StartBlock      *string
	EndBlock        *string
	ValidatorPubKey *string
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (f *SelectFilter) Filter(q *orm.Query) (*orm.Query, error) {
	if f.Addresses != nil {
		q = q.Join("LEFT OUTER JOIN transaction_outputs ON transaction_outputs.transaction_id = transaction.id").
			Join("JOIN addresses ON (addresses.id = transaction_outputs.to_address_id OR addresses.id = transaction.from_address_id)").
			WhereIn("addresses.address IN (?)", f.Addresses...)
	}

	if f.ValidatorPubKey != nil {
		q = q.Join("LEFT OUTER JOIN transaction_outputs ON transaction_outputs.transaction_id = transaction.id").
			Join("JOIN transaction_validator ON transaction_validator.transaction_id = transaction.id").
			Join("JOIN validators ON validators.public_key = ?", f.ValidatorPubKey)
	}

	if f.BlockId != nil {
		q = q.Where("transaction.block_id = ?", f.BlockId)
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

// Get transaction by hash
func (repository Repository) GetTxByHash(hash string) *models.Transaction {
	var transaction models.Transaction

	err := repository.db.Model(&transaction).Column("FromAddress").Where("hash = ?", hash).Select()
	if err != nil {
		return nil
	}

	return &transaction
}
