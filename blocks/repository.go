package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-tools/models"
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
		Column("BlockValidators", "BlockValidators.Validator.public_key").
		Apply(pagination.Filter).
		Order("id DESC").
		SelectAndCount()

	helpers.CheckErr(err)

	return blocks
}

// Get last block
func (repository Repository) GetLastBlock() models.Block {
	var block models.Block

	err := repository.DB.Model(&block).Last()
	helpers.CheckErr(err)

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

// Get slow blocks count by last 24 hours
func (repository Repository) GetSlowBlocksCountBy24h() int {
	var block models.Block

	count, err := repository.DB.Model(&block).
		Where("block_time >= ?", helpers.Seconds2Nano(6)).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).
		Count()

	helpers.CheckErr(err)

	return count
}
