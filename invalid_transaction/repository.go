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
func (r Repository) GetTxByHash(hash string) *models.InvalidTransaction {
	var transaction models.InvalidTransaction

	err := r.db.Model(&transaction).
		Relation("FromAddress").
		Relation("GasCoin").
		Where("hash = ?", hash).
		Select()

	if err != nil {
		return nil
	}

	return &transaction
}

// GetPaginatedByAddress get invalid transactions by address
func (r Repository) GetPaginatedByAddress(address string, pagination *tools.Pagination) (txs []*models.InvalidTransaction, err error) {
	pagination.Total, err = r.db.Model(&txs).
		Relation("FromAddress").
		Relation("GasCoin").
		Join(`INNER JOIN addresses ON addresses.id = "invalid_transaction"."from_address_id"`).
		Where("addresses.address = ?", address).
		Where("gas_coin_id is not null").
		Apply(pagination.Filter).
		Order("invalid_transaction.id DESC").
		SelectAndCount()

	if err != nil {
		return nil, err
	}

	return txs, nil
}
