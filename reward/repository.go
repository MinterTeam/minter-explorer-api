package reward

import (
	"github.com/MinterTeam/minter-explorer-api/v2/aggregated_reward"
	"github.com/MinterTeam/minter-explorer-api/v2/helpers"
	"github.com/MinterTeam/minter-explorer-api/v2/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg/v10"
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

type ChartData struct {
	Time   time.Time `json:"time"`
	Amount string    `json:"amount"`
}

func (repository Repository) GetAggregatedChartData(filter aggregated_reward.SelectFilter) []ChartData {
	var rewards models.AggregatedReward
	var chartData []ChartData

	err := repository.db.Model(&rewards).
		Relation("Address._").
		ColumnExpr("date_trunc('day', time_id) as time").
		ColumnExpr("SUM(amount) as amount").
		Group("time").
		Order("time").
		Apply(filter.Filter).
		Select(&chartData)

	helpers.CheckErr(err)

	return chartData
}

func (repository Repository) GetPaginatedAggregatedByAddress(filter aggregated_reward.SelectFilter, pagination *tools.Pagination) []models.AggregatedReward {
	var rewards []models.AggregatedReward
	var err error

	// get rewards
	pagination.Total, err = repository.db.Model(&rewards).
		Relation("Address.address").
		Relation("Validator").
		Apply(filter.Filter).
		Apply(pagination.Filter).
		Order("time_id DESC").
		Order("amount").
		SelectAndCount()

	helpers.CheckErr(err)

	return rewards
}
