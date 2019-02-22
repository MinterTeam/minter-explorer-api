package reward

import (
	"github.com/MinterTeam/minter-explorer-api/events"
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

// Get filtered list of rewards by Minter address
func (repository Repository) GetPaginatedByAddress(filter events.SelectFilter, pagination *tools.Pagination) []models.Reward {
	var rewards []models.Reward
	var err error

	pagination.Total, err = repository.db.Model(&rewards).
		Column("Address.address", "Validator.public_key", "Block.created_at").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("block_id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return rewards
}

type ChartData struct {
	Time   string
	Amount string
}

func (repository Repository) GetChartData(address string, scale string, startTime string, endTime string) []ChartData {
	var rewards models.Reward
	var chartData []ChartData

	err := repository.db.Model(&rewards).
		Column("Address._", "Block._").
		ColumnExpr("date_trunc(?, block.created_at) as time", scale).
		ColumnExpr("SUM(amount) as amount").
		Where("address.address = ?", address).
		Where("block.created_at >= ?", startTime).
		Where("block.created_at <= ?", endTime).
		Group("time").
		Order("time").
		Select(&chartData)

	helpers.CheckErr(err)

	return chartData
}
