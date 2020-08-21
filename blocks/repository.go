package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/core/config"
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/v2/models"
	"github.com/go-pg/pg"
	"time"
)

type Repository struct {
	DB *pg.DB
}

func NewRepository(db *pg.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

// Get block by height (id)
func (repository Repository) GetById(id uint64) *models.Block {
	var block models.Block

	err := repository.DB.Model(&block).
		Column("BlockValidators", "BlockValidators.Validator").
		Where("block.id = ?", id).
		Select()

	if err != nil {
		return nil
	}

	return &block
}

// Get paginated list of blocks
func (repository Repository) GetPaginated(pagination *tools.Pagination) []models.Block {
	var blocks []models.Block
	var err error

	pagination.Total, err = repository.DB.Model(&blocks).
		Column("BlockValidators", "BlockValidators.Validator").
		Apply(pagination.Filter).
		Order("id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return blocks
}

// Get last block
func (repository Repository) GetLastBlock() models.Block {
	var block models.Block

	repository.DB.Model(&block).Last()
	//helpers.CheckErr(err)

	return block
}

// Get average block time
func (repository Repository) GetAverageBlockTime() float64 {
	var block models.Block
	var blockTime float64

	err := repository.DB.Model(&block).
		ColumnExpr("AVG(block_time) / ?", time.Second).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&blockTime)

	helpers.CheckErr(err)

	return blockTime
}

// Get sum of delta slow time
func (repository Repository) GetSumSlowBlocksTimeBy24h() float64 {
	var block models.Block
	var sum float64

	err := repository.DB.Model(&block).
		ColumnExpr("SUM(block_time - ?) / ?", helpers.Seconds2Nano(config.SlowBlocksMaxTimeInSec), helpers.Seconds2Nano(1)).
		Where("block_time >= ?", helpers.Seconds2Nano(config.SlowBlocksMaxTimeInSec)).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Select(&sum)

	helpers.CheckErr(err)

	return sum
}
