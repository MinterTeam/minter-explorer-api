package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"time"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	orm.RegisterTable((*models.TransactionValidator)(nil))

	return &Repository{
		db: db,
	}
}

// Get paginated list of transactions by address filter
func (repository Repository) GetPaginatedTxsByAddresses(addresses []string, filter SelectFilter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Join("INNER JOIN index_transaction_by_address AS ind").
		JoinOn("ind.transaction_id = transaction.id").
		Join("INNER JOIN addresses AS a").
		JoinOn("a.id = ind.address_id").
		ColumnExpr("DISTINCT transaction.id").
		Relation("FromAddress.address").
		Relation("GasCoin").
		Column("transaction.*").
		Where("a.address IN (?)", pg.In(addresses)).
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("transaction.id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return transactions
}

// Get paginated list of transactions by select filter
func (repository Repository) GetPaginatedTxsByFilter(filter tools.Filter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		Relation("FromAddress.address").
		Relation("GasCoin").
		Column("transaction.*").
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

	err := repository.db.Model(&transaction).
		Relation("FromAddress").
		Relation("GasCoin").
		Where("hash = ?", hash).
		Select()

	if err != nil {
		return nil
	}

	return &transaction
}

type TxCountChartData struct {
	Time  time.Time
	Count uint64
}

// Get list of transactions counts filtered by created_at
func (repository Repository) GetTxCountChartDataByFilter(filter tools.Filter) []TxCountChartData {
	var tx models.Transaction
	var data []TxCountChartData

	err := repository.db.Model(&tx).
		ColumnExpr("COUNT(*) as count").
		Apply(filter.Filter).
		Select(&data)

	helpers.CheckErr(err)

	return data
}

// Get total transaction count
func (repository Repository) GetTotalTransactionCount(startTime *string) int {
	var tx models.Transaction

	query := repository.db.Model(&tx)
	if startTime != nil {
		query = query.Relation("Block._").Where("block.created_at >= ?", *startTime)
	}

	count, err := query.Count()
	helpers.CheckErr(err)

	return count
}

type Tx24hData struct {
	FeeSum float64
	Count  int
	FeeAvg float64
}

// Get transactions data by last 24 hours
func (repository Repository) Get24hTransactionsData() Tx24hData {
	var tx models.Transaction
	var data Tx24hData

	err := repository.db.Model(&tx).
		Relation("Block._").
		ColumnExpr("COUNT(*) as count, SUM(gas * gas_price) as fee_sum, AVG(gas * gas_price) as fee_avg").
		Where("block.created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&data)

	helpers.CheckErr(err)

	return data
}
