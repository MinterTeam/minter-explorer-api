package blocks

import (
	"github.com/MinterTeam/minter-explorer-api/helpers"
	"github.com/MinterTeam/minter-explorer-api/tools"
	"github.com/MinterTeam/minter-explorer-extender/models"
	"github.com/go-pg/pg"
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

	err := repository.DB.Model(&block).Column("Validators").Where("ID = ?", id).Select()
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
		Column("Validators").
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
