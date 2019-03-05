package transaction

import (
	"github.com/MinterTeam/minter-explorer-api/blocks"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"time"
)

type Repository struct {
	db *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Get paginated list of transactions by address filter
func (repository Repository) GetPaginatedTxsByAddresses(addresses []string, filter blocks.RangeSelectFilter, pagination *tools.Pagination) []models.Transaction {
	var transactions []models.Transaction
	var err error

	pagination.Total, err = repository.db.Model(&transactions).
		ColumnExpr("DISTINCT transaction.id").
		Column("transaction.*", "FromAddress.address", "TxOutput._", "TxOutput.ToAddress._").
		WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			return q.WhereIn("from_address.address IN (?)", pg.In(addresses)).
				WhereOr("tx_output__to_address.address IN (?)", pg.In(addresses)), nil
		}).
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
		Column("transaction.*", "FromAddress.address").
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
		query = query.Column("Block._").Where("block.created_at >= ?", *startTime)
	}

	count, err := query.Count()
	helpers.CheckErr(err)

	return count
}

type Tx24hData struct {
	FeeSum uint64
	Count  int
	FeeAvg float64
}

// Get transactions data by last 24 hours
func (repository Repository) Get24hTransactionsData() Tx24hData {
	var tx models.Transaction
	var data Tx24hData

	err := repository.db.Model(&tx).
		Column("Block._").
		ColumnExpr("COUNT(*) as count, SUM(gas) as fee_sum, AVG(gas) as fee_avg").
		Where("block.created_at >= ?", time.Now().AddDate(0, 0, -1).Format("2006-01-02 15:04:05")).
		Select(&data)

	helpers.CheckErr(err)

	return data
}
