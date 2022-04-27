package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/v2/core/config"
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
	orm.RegisterTable((*models.BlockValidator)(nil))

	return &Repository{
		db: db,
	}
}

// Get block by height (id)
func (r Repository) GetById(id uint64) *models.Block {
	var block models.Block

	err := r.db.Model(&block).
		Relation("BlockValidators").
		Relation("BlockValidators.Validator").
		Where("block.id = ?", id).
		Select()

	if err != nil {
		return nil
	}

	return &block
}

// Get paginated list of blocks
func (r Repository) GetPaginated(pagination *tools.Pagination) []models.Block {
	var blocks []models.Block
	var err error

	pagination.Total, err = r.db.Model(&blocks).
		Relation("BlockValidators").
		Apply(pagination.Filter).
		Order("id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return blocks
}

// Get last block
func (r Repository) GetLastBlock() models.Block {
	var block models.Block
	r.db.Model(&block).Last()
	return block
}

// Get average block time
func (r Repository) GetAverageBlockTime() float64 {
	var block models.Block
	var blockTime float64

	err := r.db.Model(&block).
		ColumnExpr("AVG(block_time) / ?", time.Second).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&blockTime)

	helpers.CheckErr(err)

	return blockTime
}

// Get sum of delta slow time
func (r Repository) GetSumSlowBlocksTimeBy24h() float64 {
	var block models.Block
	var sum float64

	err := r.db.Model(&block).
		ColumnExpr("SUM(block_time - ?) / ?", helpers.Seconds2Nano(config.SlowBlocksMaxTimeInSec), helpers.Seconds2Nano(1)).
		Where("block_time >= ?", helpers.Seconds2Nano(config.SlowBlocksMaxTimeInSec)).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&sum)

	helpers.CheckErr(err)

	return sum
}
