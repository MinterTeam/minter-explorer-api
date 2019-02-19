package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get paginated list of transactions by select filter
func (repository Repository) GetPaginatedTxsByFilter(filter SelectFilter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Join("LEFT OUTER JOIN transaction_outputs AS tx_output").
		JoinOn("tx_output.transaction_id = transaction.id").
		Join("LEFT OUTER JOIN addresses AS tx_output__to_address").
		JoinOn("tx_output__to_address.id = tx_output.to_address_id").
		Join("LEFT OUTER JOIN coins AS tx_output__coin").
		JoinOn("tx_output__coin.id = tx_output.coin_id").
		ColumnExpr("DISTINCT tx_output.id").
		Column("transaction.*", "FromAddress.address").
		ColumnExpr("tx_output.id AS tx_output__id").
		ColumnExpr("tx_output.value AS tx_output__value").
		ColumnExpr("tx_output__to_address.address AS tx_output__to_address__address").
		ColumnExpr("tx_output__coin.symbol AS tx_output__coin__symbol").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("transaction.id DESC").
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
