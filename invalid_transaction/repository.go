package invalid_transaction

import (
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// GetTxByHash get invalid transaction by hash
func (repository Repository) GetTxByHash(hash string) *models.InvalidTransaction {
	var transaction models.InvalidTransaction

	err := repository.db.Model(&transaction).Relation("FromAddress").Where("hash = ?", hash).Select()
	if err != nil {
		return nil
	}

	return &transaction
}

// GetPaginatedByAddress get invalid transactions by address
func (repository Repository) GetPaginatedByAddress(address string, pagination *tools.Pagination) (txs []*models.InvalidTransaction, err error) {
	pagination.Total, err = repository.db.
		Model(&txs).
		Relation("FromAddress").
		Join(`INNER JOIN addresses ON addresses.id = "invalid_transaction"."from_address_id"`).
		Where("addresses.address = ?", address).
		Apply(pagination.Filter).
		SelectAndCount()

	if err != nil {
		return nil, err
	}

	return txs, nil
}
