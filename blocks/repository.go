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
func (repository *Repository) GetById(id uint64) *models.Block {
	// fetch model from the database
	var block models.Block

	// fetch block
	err := repository.DB.Model(&block).Column("Validators").Where("ID = ?", id).Select()

	if err != nil {
		return nil
	}

	return &block
}

// Get paginated list of blocks
func (repository *Repository) GetPaginated(pagination *tools.Pagination) []models.Block {
	var blocks []models.Block
	var err error

	// fetch blocks
	query := repository.DB.Model(&blocks).Column("Validators")

	// apply pagination
	pagination.Total, err = pagination.ApplyFilter(query).Order("id DESC").SelectAndCount()
	helpers.CheckErr(err)

	return blocks
}
